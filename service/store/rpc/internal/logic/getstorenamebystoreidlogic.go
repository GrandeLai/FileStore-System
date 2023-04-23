package logic

import (
	"context"

	"FileStore-System/service/store/rpc/internal/svc"
	"FileStore-System/service/store/rpc/types/store"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoreNameByStoreIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStoreNameByStoreIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoreNameByStoreIdLogic {
	return &GetStoreNameByStoreIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStoreNameByStoreIdLogic) GetStoreNameByStoreId(in *store.StoreIdRequest) (*store.StoreNameResponse, error) {
	storeInfo, err := l.svcCtx.StoreModel.FindOne(l.ctx, in.StoreId)
	if err != nil {
		return nil, err
	}
	return &store.StoreNameResponse{StoreName: storeInfo.Name}, nil
}
