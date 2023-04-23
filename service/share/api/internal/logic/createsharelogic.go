package logic

import (
	"FileStore-System/common/utils"
	"FileStore-System/service/share/api/internal/svc"
	"FileStore-System/service/share/api/internal/types"
	"FileStore-System/service/share/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type CreateShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShareLogic {
	return &CreateShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShareLogic) CreateShare(req *types.CreateShareRequest) (resp *types.CreateShareResponse, err error) {

	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	shareURL := shortuuid.New()
	var datas []*model.Share
	for _, storeId := range req.StoreIds {
		storeId, err := strconv.ParseInt(storeId, 10, 64)
		if err != nil {
			return nil, err
		}
		data := &model.Share{
			Id:          utils.GenerateNewId(l.svcCtx.RedisClient, "share"),
			UserId:      userId,
			StoreId:     storeId,
			ExpiredTime: req.ExpiredTime,
			ShareUrl:    shareURL,
		}
		datas = append(datas, data)
	}
	_, err = l.svcCtx.ShareModel.InsertByBatch(l.ctx, datas)
	if err != nil {
		return nil, err
	}
	return &types.CreateShareResponse{ShareURL: shareURL}, nil
}
