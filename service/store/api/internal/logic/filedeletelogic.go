package logic

import (
	"FileStore-System/common/errorx"
	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"
	"FileStore-System/service/user/rpc/types/user"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDeleteLogic) FileDelete(req *types.StoreDeleteRequest) (resp *types.StoreDeleteResponse, err error) {
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}

	parentId, err := strconv.ParseInt(req.ParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	userFile, err := l.svcCtx.UserStoreModel.FindByFactors(l.ctx, parentId, userId, id)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.UserStoreModel.Delete(l.ctx, userFile.Id)
	if err != nil {
		return nil, errorx.NewDefaultError("更新个人存储池失败！")
	}

	store, err := l.svcCtx.StoreModel.FindOne(l.ctx, userFile.StoreId)
	if err != nil {
		return nil, errorx.NewDefaultError("中心存储池找不到该数据！")
	}

	//查询是否还有人关联这个存储文件
	count, err := l.svcCtx.UserStoreModel.CountStoreUsage(l.ctx, userFile.Id)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		err := l.svcCtx.StoreModel.Delete(l.ctx, userFile.Id)
		if err != nil {
			return nil, errorx.NewDefaultError("更新中心存储池失败！")
		}
		//TODO：删除存储的文件
		path := strings.Split(store.Path, "/")
		objectName := path[len(path)-1]
		err = l.svcCtx.MinioClient.RemoveObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError("服务器删除对象失败")
		}
	}

	_, err = l.svcCtx.UserRpc.DecreaseVolume(l.ctx, &user.DecreaseVolumeReq{
		Id:   userId,
		Size: store.Size,
	})
	if err != nil {
		return nil, errorx.NewDefaultError("更新容量失败！")
	}
	return &types.StoreDeleteResponse{}, nil
}
