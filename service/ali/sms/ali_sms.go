package aliyum

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type SmsService struct {
	AliAccessKey    string
	AliAccessSecret string
	SignName        string
}

type AliRe struct {
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	BizId     string `json:"BizId,omitempty"`
}

func NewSmsService(aliAccessKey, aliAccessSecret, signName string) *SmsService {
	return &SmsService{
		AliAccessKey:    aliAccessKey,
		AliAccessSecret: aliAccessSecret,
		SignName:        signName,
	}
}

func (s SmsService) SmsSend(phoneNum, templateCode string, templateParam map[string]string) (bool, error) {
	if len(s.AliAccessKey) == 0 || len(s.AliAccessSecret) == 0 {
		return false, errors.New("配置加载失败")
	}
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", s.AliAccessKey, s.AliAccessSecret)
	if err != nil {
		return false, err
	}
	templateParamStr, _ := json.Marshal(templateParam)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phoneNum
	request.QueryParams["SignName"] = s.SignName
	request.QueryParams["TemplateCode"] = templateCode              //短信发送模版
	request.QueryParams["TemplateParam"] = string(templateParamStr) //短信发送模版变量

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return false, err
	}
	re := AliRe{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &re)
	if err != nil {
		return false, err
	}
	if re.Code == "OK" && re.Message == "OK" {
		return true, nil
	}
	return false, errors.New(re.Message)
}

func (s SmsService) SmsCode(smsCodeTemplate, phoneNum, code string) (bool, error) {
	templateParam := map[string]string{
		"code": code,
	}
	return s.SmsSend(phoneNum, smsCodeTemplate, templateParam)
}
