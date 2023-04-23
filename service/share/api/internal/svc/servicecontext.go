package svc

import (
	"FileStore-System/service/share/api/internal/config"
	"FileStore-System/service/share/model"
	"FileStore-System/service/store/rpc/storeclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	ShareModel  model.ShareModel
	StoreRpc    storeclient.Store
	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		ShareModel: model.NewShareModel(conn, c.CacheRedis),
		StoreRpc:   storeclient.NewStore(zrpc.MustNewClient(c.StoreRpc)),
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
	}
}
