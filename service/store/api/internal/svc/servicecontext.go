package svc

import (
	"FileStore-System/service/store/api/internal/config"
	minio2 "FileStore-System/service/store/api/internal/svc/minioclient"
	"FileStore-System/service/store/model"
	"FileStore-System/service/user/rpc/userclient"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	StoreModel     model.StoreModel
	UserStoreModel model.UserStoreModel
	UserRpc        userclient.User
	MinioClient    *minio.Client
	RedisClient    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		StoreModel:     model.NewStoreModel(conn, c.CacheRedis),
		UserStoreModel: model.NewUserStoreModel(conn, c.CacheRedis),
		UserRpc:        userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		MinioClient: minio2.NewMinioClient(c),
	}
}
