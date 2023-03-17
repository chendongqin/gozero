package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"
)

// 用户上传文件时指定的前缀。
// var uploadDir string = "user-dir-prefix/"
var expireTime int64 = 30

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type OssConfig struct {
	OssUrl          string // oss url路径
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string // oss图片endpoint
	Bucket          string // oss图片bucket
	CallbackUrl     string // oss图片回调地址
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

type OssSvr struct {
	CallbackUrl     string
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
}

func NewOssSvr(callbackUrl, accessKeyId, accessKeySecret, endpoint string) *OssSvr {
	return &OssSvr{
		CallbackUrl:     callbackUrl,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Endpoint:        endpoint,
	}
}

func (s *OssSvr) GetPolicyTokenById(uploadDir, tmpKey, upType string, id int64) *PolicyToken {
	callBackBody := fmt.Sprintf("filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}&tmp_key=%s&class=%s&uid=%d", tmpKey, upType, id)
	callBackUrl := s.CallbackUrl
	// fmt.Println(avatarCallBackUrl)
	return s.commUploadPolicyToken(uploadDir, callBackBody, callBackUrl)
}

func (s *OssSvr) commUploadPolicyToken(uploadDir, callBackBody, callBackUrlString string) *PolicyToken {
	now := time.Now().Unix()
	expireEnd := now + expireTime
	var tokenExpire = getGmtIso8601(expireEnd)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, _ := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(s.AccessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callBackUrlString
	callbackParam.CallbackBody = callBackBody
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callbackStr, _ := json.Marshal(callbackParam)

	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	policyToken := &PolicyToken{
		AccessKeyId: s.AccessKeyId,
		Host:        s.Endpoint,
		Expire:      expireEnd,
		Signature:   signedStr,
		Directory:   uploadDir,
		Policy:      debyte,
		Callback:    callbackBase64,
	}

	return policyToken
}
