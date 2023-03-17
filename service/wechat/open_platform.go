package wechat

import (
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/silenceper/wechat/v2/openplatform"
)

type OpenPlatformService struct {
	Service *openplatform.OpenPlatform
	AppId   string
}

//NewOpenPlatformService 初始化开放平台服务
func NewOpenPlatformService(wechatConf *WechatConfig, redisConf *WechatRedisConfig, appid string) *OpenPlatformService {
	openPlatform := InitOpenPlatformService(wechatConf, redisConf)
	return &OpenPlatformService{
		Service: openPlatform,
		AppId:   appid,
	}
}

func (s *OpenPlatformService) GetUserAccessToken(code string) (result oauth.ResAccessToken, err error) {
	return s.Service.GetOfficialAccount(s.AppId).GetOauth().GetUserAccessToken(code)
}

func (s *OpenPlatformService) GetUserInfoByCode(code string) (result oauth.UserInfo, err error) {
	res, err := s.GetUserAccessToken(code)
	if err != nil {
		return oauth.UserInfo{}, err
	}
	return s.GetUserInfo(res.AccessToken, res.OpenID, "")
}

func (s *OpenPlatformService) GetUserInfo(accessToken, openID, lang string) (result oauth.UserInfo, err error) {
	return s.Service.GetOfficialAccount(s.AppId).GetOauth().GetUserInfo(accessToken, openID, lang)
}
