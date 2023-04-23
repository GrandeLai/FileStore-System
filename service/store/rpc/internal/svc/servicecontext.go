package svc

import (
	"FileStore-System/service/store/model"
	"FileStore-System/service/store/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	StoreModel     model.StoreModel
	UserStoreModel model.UserStoreModel
	RedisClient    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		StoreModel:     model.NewStoreModel(conn, c.CacheRedis),
		UserStoreModel: model.NewUserStoreModel(conn, c.CacheRedis),
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
	}
}
