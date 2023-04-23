package logic

import (
	"FileStore-System/common/utils"
	"FileStore-System/service/store/model"
	"FileStore-System/service/store/rpc/internal/svc"
	"FileStore-System/service/store/rpc/types/store"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateByShareInBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateByShareInBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateByShareInBatchLogic {
	return &CreateByShareInBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateByShareInBatchLogic) CreateByShareInBatch(in *store.CreateByShareInBatchRequest) (*store.CreateByShareInBatchResponse, error) {
	var userStores []*model.UserStore
	for _, req := range in.CreateByShareRequest {
		newId := utils.GenerateNewId(l.svcCtx.RedisClient, "user_store")
		userStore := &model.UserStore{
			Id:       newId,
			UserId:   req.UserId,
			ParentId: req.ParentId,
			StoreId:  req.StoreId,
			Name:     req.Name,
		}
		userStores = append(userStores, userStore)
	}
	_, err := l.svcCtx.UserStoreModel.InsertByBatch(l.ctx, userStores)
	if err != nil {
		return nil, err
	}
	return &store.CreateByShareInBatchResponse{}, nil
}
