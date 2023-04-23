package logic

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"io/ioutil"
	"net/http"
	"strconv"

	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	writer http.ResponseWriter
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer http.ResponseWriter) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		writer: writer,
	}
}

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadRequest) error {
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return err
	}
	file, err := l.svcCtx.StoreModel.FindOne(l.ctx, id)
	if err != nil {
		return errors.New("未找到文件")
	}
	object, err := l.svcCtx.MinioClient.GetObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, file.Path, minio.GetObjectOptions{})
	if err != nil {
		return errors.New("获取文件失败")
	}
	data, err := ioutil.ReadAll(object)
	if err != nil {
		return errors.New("获取文件失败")
	}

	l.writer.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	l.writer.Header().Set("content-disposition", "attachment; filename=\""+file.Name+"\"")
	l.writer.Write(data)
	return nil
}
