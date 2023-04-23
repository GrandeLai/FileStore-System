package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileNameUpdateLogic {
	return &FileNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileNameUpdateLogic) FileNameUpdate(req *types.FileNameUpdateRequest) (resp *types.FileNameUpdateResponse, err error) {
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	parentId, err := strconv.ParseInt(req.ParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	userFileInfo, err := l.svcCtx.UserStoreModel.FindByFactors(l.ctx, parentId, userId, id)
	if err != nil {
		return nil, errors.New("无法在此文件夹下找到该文件")
	}
	userFileInfo.Name = req.Name
	err = l.svcCtx.UserStoreModel.Update(l.ctx, userFileInfo)
	if err != nil {
		return nil, err
	}
	return &types.FileNameUpdateResponse{}, nil
}
