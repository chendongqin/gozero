package svc

import (
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/rest"
	"gozore-mall/app/cmd/admin/internal/config"
	"gozore-mall/app/cmd/admin/internal/middleware"
	"gozore-mall/service/gormpool"
)

type ServiceContext struct {
	Config       config.Config
	Redis        *redis.Redis
	PermMenuAuth rest.Middleware
	AdminActLog  rest.Middleware
	GormQuery    *gormpool.CommonQuery
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
		cache.NewStat("zhg-admin"), ErrNoRows)
	signalChan := make(gormpool.CacheSignalChan, 1)
	commonQuery := gormpool.NewCommonQuery(db, writeDb, redisCache, &gormpool.CacheSignal{
		CacheOff:     false,
		CacheErr:     nil,
		CacheErrFile: c.CacheErrFile,
		CacheChan:    signalChan,
	})
	return &ServiceContext{
		Config:       c,
		Redis:        redisClient,
		PermMenuAuth: middleware.NewPermMenuAuthMiddleware(redisClient, commonQuery).Handle,
		AdminActLog:  middleware.NewAdminActLogMiddleware(commonQuery).Handle,
		GormQuery:    commonQuery,
	}
}
