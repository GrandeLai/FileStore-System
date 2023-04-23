package logic

import (
	"FileStore-System/common/errorx"
	"FileStore-System/common/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileListLogic) FileList(req *types.FileListRequest) (resp *types.FileListResponse, err error) {
	//获得分页的初始下标和每页大小
	pageSize := req.Size
	if req.Size == 0 {
		pageSize = utils.DefaultPageSize
	}
	startPage := req.Page
	if startPage == 0 {
		startPage = 1
	}
	startIndex := pageSize * (startPage - 1)
	//根据文件夹id，然后作为父目录id去搜目录下的数据
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	parentId, err := strconv.ParseInt(req.ParentId, 10, 64)
	if err != nil {
		return nil, err
	}
	folderInfo, err := l.svcCtx.StoreModel.FindOne(l.ctx, parentId)
	if err != nil {
		return nil, err
	}
	if folderInfo.IsFolder != 1 {
		return nil, errors.New("不是文件夹")
	}
	allUserStore, err := l.svcCtx.UserStoreModel.FindAllInPage(l.ctx, parentId, userId, startIndex, pageSize)
	if err != nil {
		return nil, errorx.NewDefaultError("该文件夹下搜索文件失败！")
	}
	newList := make([]*types.File, 0)
	for _, userStore := range allUserStore {
		storeInfo, err := l.svcCtx.StoreModel.FindOne(l.ctx, userStore.StoreId)
		if err != nil {
			return nil, err
		}
		newList = append(newList, &types.File{
			Id:   string([]rune(strconv.FormatInt(userStore.Id, 10))),
			Name: storeInfo.Name,
			Ext:  storeInfo.Ext,
			Path: storeInfo.Path,
			Size: storeInfo.Size,
		})
	}
	return &types.FileListResponse{
		List:  newList,
		Count: int64(len(allUserStore)),
	}, err
}
