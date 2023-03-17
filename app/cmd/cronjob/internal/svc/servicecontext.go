package svc

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
	"gozore-api/service/gormpool"
	"gozore-mall/app/cmd/cronjob/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *redis.Redis
	GormQuery *gormpool.CommonQuery
}

var ErrNoRows = errors.New("server cache: no rows in result set")

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gormpool.CoonMysql(c.GormConfig)
	if err != nil {
		panic(err)
	}
	writeDb, err := gormpool.CoonMysql(c.WriteGormConfig)
	if err != nil {
		panic(err)
	}
	redisClient := redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = c.Redis.Type
		r.Pass = c.Redis.Pass
	})
	redisCache := cache.New(c.CacheRedisCluster, syncx.NewSingleFlight(),
		cache.NewStat("ad-skin-api"), ErrNoRows)
	signalChan := make(gormpool.CacheSignalChan, 1)
	commonQuery := gormpool.NewCommonQuery(db, writeDb, redisCache, &gormpool.CacheSignal{
		CacheOff:     false,
		CacheErr:     nil,
		CacheErrFile: c.CacheErrFile,
		CacheChan:    signalChan,
	})
	return &ServiceContext{
		Config:    c,
		Redis:     redisClient,
		GormQuery: commonQuery,
	}
}
