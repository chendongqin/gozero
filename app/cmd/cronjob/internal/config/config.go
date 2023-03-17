package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"gozore-api/service/gormpool"
)

type Config struct {
	service.ServiceConf
	UserCenterRpcConf zrpc.RpcClientConf
	GormConfig        gormpool.GormConfig
	WriteGormConfig   gormpool.GormConfig
	CacheRedisCluster cache.CacheConf
	CacheErrFile      string
	Redis             struct {
		Host string
		Pass string
		Type string
	}

	Alipay struct {
		AppId        string
		PrivateKey   string
		CertPathApp  string
		CertPathRoot string
		CertPathAli  string
		NotifyUrl    string
		ReturnUrl    string
	}
}
