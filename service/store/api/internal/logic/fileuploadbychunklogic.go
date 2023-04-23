package logic

import (
	"context"
	"mime/multipart"

	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadByChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadByChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadByChunkLogic {
	return &FileUploadByChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadByChunkLogic) FileUploadByChunk(req *types.FileUploadByChunkRequest, file multipart.File, fileHeader *multipart.FileHeader) (resp *types.FileUploadByChunkResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
