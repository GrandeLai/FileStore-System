// Code generated by goctl. DO NOT EDIT!
// Source: store.proto

package server

import (
	"context"

	"FileStore-System/service/store/rpc/internal/logic"
	"FileStore-System/service/store/rpc/internal/svc"
	"FileStore-System/service/store/rpc/types/store"
)

type StoreServer struct {
	svcCtx *svc.ServiceContext
	store.UnimplementedStoreServer
}

func NewStoreServer(svcCtx *svc.ServiceContext) *StoreServer {
	return &StoreServer{
		svcCtx: svcCtx,
	}
}

func (s *StoreServer) GetStoreNameByStoreId(ctx context.Context, in *store.StoreIdRequest) (*store.StoreNameResponse, error) {
	l := logic.NewGetStoreNameByStoreIdLogic(ctx, s.svcCtx)
	return l.GetStoreNameByStoreId(in)
}

func (s *StoreServer) CreateByShare(ctx context.Context, in *store.CreateByShareRequest) (*store.CreateByShareResponse, error) {
	l := logic.NewCreateByShareLogic(ctx, s.svcCtx)
	return l.CreateByShare(in)
}

func (s *StoreServer) CreateByShareInBatch(ctx context.Context, in *store.CreateByShareInBatchRequest) (*store.CreateByShareInBatchResponse, error) {
	l := logic.NewCreateByShareInBatchLogic(ctx, s.svcCtx)
	return l.CreateByShareInBatch(in)
}

func (s *StoreServer) GetStoreInfoInBatchByStoreId(ctx context.Context, in *store.StoreIdsRequest) (*store.StoreInfosResponse, error) {
	l := logic.NewGetStoreInfoInBatchByStoreIdLogic(ctx, s.svcCtx)
	return l.GetStoreInfoInBatchByStoreId(in)
}