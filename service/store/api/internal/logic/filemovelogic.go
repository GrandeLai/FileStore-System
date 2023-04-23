package logic

import (
	"FileStore-System/common/errorx"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileMoveLogic {
	return &FileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileMoveLogic) FileMove(req *types.FileMoveRequest) (resp *types.FileMoveResponse, err error) {
	parentId, err := strconv.ParseInt(req.ParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	targetParentId, err := strconv.ParseInt(req.TargetParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	//检测该文件是否存在
	userFileInfo, err := l.svcCtx.UserStoreModel.FindByFactors(l.ctx, parentId, userId, id)
	if err != nil {
		return nil, errorx.NewDefaultError("原文件不存在！")
	}
	//检测新目录是否已存在该文件
	count, err := l.svcCtx.UserStoreModel.CountByIdAndParentId(l.ctx, id, userId, targetParentId)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errorx.NewDefaultError("已存在相同名称的文件！")
	}
	//修改
	userFileInfo.ParentId = targetParentId
	err = l.svcCtx.UserStoreModel.Update(l.ctx, userFileInfo)
	if err != nil {
		return nil, err
	}
	return &types.FileMoveResponse{}, nil
}
