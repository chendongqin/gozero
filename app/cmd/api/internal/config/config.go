package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"gozore-mall/service/gormpool"
)

type Config struct {
	rest.RestConf
	Salt    string
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	GormConfig        gormpool.GormConfig
	WriteGormConfig   gormpool.GormConfig
	CacheErrFile      string
	CacheRedisCluster cache.CacheConf
	Redis             struct {
		Host string
		Pass string
		Type string
	}
	Aliyum struct {
		AccessKey       string
		AccessSecret    string
		SignName        string
		SmsCodeTemplate string
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
