package logic

import (
	"FileStore-System/common/errorx"
	"FileStore-System/service/user/rpc/types/user"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FolderDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFolderDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FolderDeleteLogic {
	return &FolderDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FolderDeleteLogic) FolderDelete(req *types.StoreDeleteRequest) (resp *types.StoreDeleteResponse, err error) {
	parentId, err := strconv.ParseInt(req.ParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
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
	filesInFolder, err := l.svcCtx.UserStoreModel.FindAllInPage(l.ctx, id, userId, 0, 0)
	if err != nil {
		return nil, err
	}
	if len(filesInFolder) > 0 {
		err := l.svcCtx.UserStoreModel.DeleteByParentId(l.ctx, id)
		if err != nil {
			return nil, errorx.NewDefaultError("删除文件夹下文件失败！")
		}
		var sumVolume int64
		for _, files := range filesInFolder {
			store, err := l.svcCtx.StoreModel.FindOne(l.ctx, files.StoreId)
			if err != nil {
				return nil, errorx.NewDefaultError("服务器错误！")
			}
			sumVolume += store.Size
		}
		_, err = l.svcCtx.UserRpc.DecreaseVolume(l.ctx, &user.DecreaseVolumeReq{
			Id:   userId,
			Size: sumVolume,
		})
		if err != nil {
			return nil, errorx.NewDefaultError("更新容量失败！")
		}
	}
	return &types.StoreDeleteResponse{}, nil
}
