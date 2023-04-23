package logic

import (
	"FileStore-System/service/store/rpc/types/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"FileStore-System/service/share/api/internal/svc"
	"FileStore-System/service/share/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveFromShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveFromShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveFromShareLogic {
	return &SaveFromShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveFromShareLogic) SaveFromShare(req *types.SaveFromShareRequest) (resp *types.SaveFromShareResponse, err error) {
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	toParentId, err := strconv.ParseInt(req.ToParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	shares, err := l.svcCtx.ShareModel.FindByShareURL(l.ctx, req.ShareURL)
	if err != nil {
		return nil, err
	}
	//判断保存的文件是否全部在shareurl的文件列表中
	if len(shares) < len(req.StoreIds) {
		return nil, errors.New("保存文件数量不符")
	}
	m := make(map[string]bool)
	for _, share := range shares {
		m[string([]rune(strconv.FormatInt(share.StoreId, 10)))] = true
	}
	for _, storeId := range req.StoreIds {
		if !m[storeId] {
			return nil, errors.New("无法保存所有文件")
		}
	}
	if len(req.StoreIds) == 1 {
		storeId, err := strconv.ParseInt(req.StoreIds[0], 10, 64)
		if err != nil {
			return nil, err
		}
		_, err = l.svcCtx.StoreRpc.CreateByShare(l.ctx, &store.CreateByShareRequest{
			UserId:   userId,
			ParentId: toParentId,
			StoreId:  storeId,
		})
		if err != nil {
			return nil, err
		}
		return &types.SaveFromShareResponse{}, nil
	} else {
		var batches []*store.CreateByShareRequest
		for _, storeId := range req.StoreIds {
			storeId, err := strconv.ParseInt(storeId, 10, 64)
			if err != nil {
				return nil, err
			}
			batch := &store.CreateByShareRequest{
				UserId:   userId,
				ParentId: toParentId,
				StoreId:  storeId,
			}
			batches = append(batches, batch)
		}
		_, err := l.svcCtx.StoreRpc.CreateByShareInBatch(l.ctx, &store.CreateByShareInBatchRequest{CreateByShareRequest: batches})
		if err != nil {
			return nil, err
		}
		return &types.SaveFromShareResponse{}, nil
	}
}
