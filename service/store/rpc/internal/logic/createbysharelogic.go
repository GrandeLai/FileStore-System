package logic

import (
	"FileStore-System/common/utils"
	"FileStore-System/service/store/model"
	"context"

	"FileStore-System/service/store/rpc/internal/svc"
	"FileStore-System/service/store/rpc/types/store"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateByShareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateByShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateByShareLogic {
	return &CreateByShareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateByShareLogic) CreateByShare(in *store.CreateByShareRequest) (*store.CreateByShareResponse, error) {
	newId := utils.GenerateNewId(l.svcCtx.RedisClient, "user_repository")
	_, err := l.svcCtx.UserStoreModel.Insert(l.ctx, &model.UserStore{
		Id:       newId,
		UserId:   in.UserId,
		ParentId: in.ParentId,
		StoreId:  in.StoreId,
		Name:     in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &store.CreateByShareResponse{Id: newId}, nil
}
