package wechat

import (
	"context"
	wechatV2 "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/openplatform"
	openplatformConfig "github.com/silenceper/wechat/v2/openplatform/config"
)

type WechatRedisConfig struct {
	Host     string
	Database int
	Password string
}

type WechatConfig struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
}

func InitOpenPlatformService(wechatConf *WechatConfig, redisConf *WechatRedisConfig) *openplatform.OpenPlatform {
	wc := wechatV2.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        redisConf.Host,
		Database:    redisConf.Database,
		Password:    redisConf.Password,
		MaxActive:   10,
		MaxIdle:     10,
		IdleTimeout: 60, //second
	}
	redisCache := cache.NewRedis(context.Background(), redisOpts)
	cfg := &openplatformConfig.Config{
		AppID:          wechatConf.AppID,
		AppSecret:      wechatConf.AppSecret,
		Token:          wechatConf.Token,
		EncodingAESKey: wechatConf.EncodingAESKey,
		Cache:          redisCache,
	}
	openPlatform := wc.GetOpenPlatform(cfg)
	return openPlatform
}
