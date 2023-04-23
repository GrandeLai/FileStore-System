package logic

import (
	"FileStore-System/service/store/model"
	"FileStore-System/service/store/rpc/internal/svc"
	"FileStore-System/service/store/rpc/types/store"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoreInfoInBatchByStoreIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStoreInfoInBatchByStoreIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoreInfoInBatchByStoreIdLogic {
	return &GetStoreInfoInBatchByStoreIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStoreInfoInBatchByStoreIdLogic) GetStoreInfoInBatchByStoreId(in *store.StoreIdsRequest) (*store.StoreInfosResponse, error) {
	var stores []*model.Store
	for _, storeId := range in.StoreIds {
		store, err := l.svcCtx.StoreModel.FindOne(l.ctx, storeId)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	var storeInfos []*store.StoreInfo
	for _, storeI := range stores {
		storeInfo := &store.StoreInfo{
			StoreId: storeI.Id,
			Ext:     storeI.Ext,
			Size:    storeI.Size,
			Path:    storeI.Path,
			Name:    storeI.Name,
		}
		storeInfos = append(storeInfos, storeInfo)
	}
	return &store.StoreInfosResponse{StoreInfo: storeInfos}, nil
}
