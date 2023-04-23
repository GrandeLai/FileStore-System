package logic

import (
	"FileStore-System/common/errorx"
	"FileStore-System/common/utils"
	"FileStore-System/service/store/api/internal/svc"
	"FileStore-System/service/store/api/internal/types"
	"FileStore-System/service/store/model"
	"FileStore-System/service/user/rpc/types/user"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"mime/multipart"
	"path"
	"strconv"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest, parentId int64, file multipart.File, fileHeader *multipart.FileHeader) (resp *types.FileUploadResponse, err error) {

	// 判断是否已达用户容量上限
	userId, err := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId"))).Int64()
	if err != nil {
		return nil, err
	}
	volumeInfo, err := l.svcCtx.UserRpc.GetVolumeByUserId(l.ctx, &user.VolumeDetailReq{Id: userId})
	if err != nil {
		return nil, err
	}
	if volumeInfo.NowVolume+fileHeader.Size > volumeInfo.TotalVolume {
		return nil, errorx.NewDefaultError("文件过大！")
	}
	// 增加用户当前已存储容量
	_, err = l.svcCtx.UserRpc.AddVolume(l.ctx, &user.AddVolumeReq{
		Id:   userId,
		Size: fileHeader.Size,
	})
	if err != nil {
		return nil, err
	}
	// 判断文件是否已存在，若已存在则为秒传成功
	//b := make([]byte, fileHeader.Size)
	//_, err = b.Read(b)
	//if err != nil {
	//	return nil, err
	//}
	sha := sha256.New()
	io.Copy(sha, file)
	file.Seek(0, 0) // 将文件指针归位
	hash := hex.EncodeToString(sha.Sum(nil))
	count, err := l.svcCtx.StoreModel.CountByHash(l.ctx, hash)
	if count > 0 {
		storeInfo, err := l.svcCtx.StoreModel.FindStoreByHash(l.ctx, hash)
		if err != nil {
			return nil, err
		}
		return &types.FileUploadResponse{Id: string(storeInfo.Id)}, err
	}
	newId1 := utils.GenerateNewId(l.svcCtx.RedisClient, "store")
	//contentType := fileHeader.Header.Get("Content-Type")
	var filePath string
	if utils.ObjectStorageType == "minio" {
		objectName := utils.GenerateUUID() + path.Ext(fileHeader.Filename)
		_, err = l.svcCtx.MinioClient.PutObject(context.Background(), "filestore", objectName, file, fileHeader.Size,
			minio.PutObjectOptions{ContentType: "binary/octet-stream"})
		if err != nil {
			//return nil, errorx.NewDefaultError("服务器处理上传失败！")
			return nil, err
		}
		filePath = l.svcCtx.Config.MinIO.BucketName + "/" + objectName
	} else {

	}
	// 插入数据
	_, err = l.svcCtx.StoreModel.InsertWithId(l.ctx, &model.Store{
		Id:       newId1,
		Hash:     hash,
		Ext:      path.Ext(fileHeader.Filename),
		Size:     fileHeader.Size,
		Path:     filePath,
		Name:     fileHeader.Filename,
		IsFolder: 0,
		Status:   0,
	})
	if err != nil {
		return nil, errorx.NewDefaultError("上传失败！")
	}
	newId2 := utils.GenerateNewId(l.svcCtx.RedisClient, "user_store")
	_, err = l.svcCtx.UserStoreModel.InsertWithId(l.ctx, &model.UserStore{
		Id:       newId2,
		UserId:   userId,
		ParentId: parentId,
		StoreId:  newId1,
		Name:     fileHeader.Filename,
	})
	if err != nil {
		return nil, errorx.NewDefaultError("上传失败！")
	}
	return &types.FileUploadResponse{Id: string([]rune(strconv.FormatInt(newId1, 10)))}, err
}
