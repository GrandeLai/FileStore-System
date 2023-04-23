package logic

import (
	"FileStore-System/common/errorx"
	"FileStore-System/common/utils"
	"FileStore-System/service/share/api/internal/svc"
	"FileStore-System/service/share/api/internal/types"
	"FileStore-System/service/store/rpc/types/store"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"strconv"
)

type DetailShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailShareLogic {
	return &DetailShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailShareLogic) DetailShare(req *types.ShareDetailRequest) (resp *types.ShareDetailResponse, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//先查询缓存中有没有该数据
	redisQueryKey := utils.CacheShareKey + req.ShareURL
	ifExists, err := l.svcCtx.RedisClient.Exists(redisQueryKey)
	if err != nil {
		return nil, err
	}
	if ifExists == true {
		//有
		jsonStr, err := l.svcCtx.RedisClient.Get(redisQueryKey)
		if err != nil {
			return nil, err
		}
		//判断数据是否为空
		if jsonStr == "" {
			return nil, errorx.NewCodeError(100, "查无此分享信息")
		}
		var shareInfo types.ShareDetailResponse
		err = json.UnmarshalFromString(jsonStr, &shareInfo)
		if err != nil {
			return nil, err
		}
		return &shareInfo, nil
	}
	//从数据库查询数据
	//申请分布式锁，获取repositoryId和返回user库、repositoryPool库的相应值
	redisLockKey := redisQueryKey
	redisLock := redis.NewRedisLock(l.svcCtx.RedisClient, redisLockKey)
	redisLock.SetExpire(utils.RedisLockExpireSeconds)
	if ok, err := redisLock.Acquire(); !ok || err != nil {
		return nil, errorx.NewCodeError(100, "当前有其他用户正在进行操作，请稍后重试")
	}
	defer func() {
		recover()
		redisLock.Release()
	}()
	shareInfos, err := l.svcCtx.ShareModel.FindByShareURL(l.ctx, req.ShareURL)
	switch err {
	case nil:
		break
	case sqlc.ErrNotFound:
		//缓存空数据
		err = l.svcCtx.RedisClient.Setex(redisQueryKey, "", utils.RedisLockExpireSeconds)
		if err != nil {
			return nil, err
		}
		return nil, errorx.NewCodeError(100, "查无此分享信息")
	default:
		return nil, err
	}
	var storeIds []int64
	for _, shareInfo := range shareInfos {
		storeIds = append(storeIds, shareInfo.StoreId)
	}
	storeInfos, err := l.svcCtx.StoreRpc.GetStoreInfoInBatchByStoreId(l.ctx, &store.StoreIdsRequest{StoreIds: storeIds})
	if err != nil {
		return nil, errorx.NewDefaultError("无法获得用户储存库的信息！")
	}
	var shareDetails []*types.ShareDetail
	for i, storeInfo := range storeInfos.StoreInfo {
		shareDetail := &types.ShareDetail{
			StoreId:   string([]rune(strconv.FormatInt(storeInfo.StoreId, 10))),
			Name:      storeInfo.Name,
			Size:      storeInfo.Size,
			CreatedAt: shareInfos[i].CreateTime.String(),
			UpdatedAt: shareInfos[i].UpdateTime.String(),
		}
		shareDetails = append(shareDetails, shareDetail)
	}
	//把数据存储到缓存中
	DetailInfo := types.ShareDetailResponse{ShareDetails: shareDetails}
	jsonStr, err := json.MarshalToString(DetailInfo)
	if err != nil {
		return nil, err
	}
	l.svcCtx.RedisClient.Setex(redisQueryKey, jsonStr, utils.RedisLockExpireSeconds)
	return &DetailInfo, nil
}
