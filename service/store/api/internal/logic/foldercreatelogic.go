package logic

import (
	"FileStore-System/common/errorx"
	"FileStore-System/common/utils"
	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"
	"FileStore-System/service/store/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type FolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FolderCreateLogic {
	return &FolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FolderCreateLogic) FolderCreate(req *types.FolderCreateRequest) (resp *types.FolderCreateResponse, err error) {
	parentId, err := strconv.ParseInt(req.ParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	//验证文件夹名字不存在：
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	existCount, err := l.svcCtx.UserStoreModel.CountByParentIdAndName(l.ctx, parentId, userId, req.Name)
	if err != nil {
		return nil, errorx.NewDefaultError("服务器错误！")
	}
	if existCount > 0 {
		return nil, errorx.NewDefaultError("已存在相同名称的文件夹！")
	}
	newId1 := utils.GenerateNewId(l.svcCtx.RedisClient, "store")
	_, err = l.svcCtx.StoreModel.InsertWithId(l.ctx, &model.Store{
		Id:       newId1,
		Name:     req.Name,
		Hash:     string([]rune(strconv.FormatInt(newId1, 10))),
		Status:   0,
		IsFolder: 1,
	})
	if err != nil {
		return nil, err
	}
	newId2 := utils.GenerateNewId(l.svcCtx.RedisClient, "user_store")
	_, err = l.svcCtx.UserStoreModel.InsertWithId(l.ctx, &model.UserStore{
		Id:       newId2,
		UserId:   userId,
		ParentId: parentId,
		StoreId:  newId1,
		Name:     req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &types.FolderCreateResponse{Id: string([]rune(strconv.FormatInt(newId1, 10)))}, nil
}
