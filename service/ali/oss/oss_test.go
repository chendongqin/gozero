package oss

import (
	"fmt"
	"gozore-mall/common/utils"
	"strings"
	"testing"
	"time"
)

func TestOssToken(t *testing.T) {
	fileType := "jpg"
	userId := int64(1000)
	uploadUrl := fmt.Sprintf("%s/%s/%d-%d.%s", "user", time.Now().Format("20060102"), userId, time.Now().Unix(), fileType)
	randKey := utils.GetRandom(13, 2)
	md5Str := utils.Md5Encode(fmt.Sprintf("%s|%d|%s|imguploadsafe", randKey, userId, "user"))
	randKey = randKey + strings.ToUpper(md5Str[0:3])
	// RedisValue.SetUploadOssTmp(userId, class, randKey, uploadUrl)
	ossSvr := NewOssSvr("", "", "", "")
	policyToken := ossSvr.GetPolicyTokenById(uploadUrl, randKey, "user", userId)
	ossUrl := "oss bucket url"
	bucketName := "oss bucket name"

	res := map[string]interface{}{
		"policy_token": policyToken,
		"oss_url":      ossUrl,
		"bucket_name":  bucketName,
	}
	fmt.Println(res)

	return
}

func TestCheckBack(t *testing.T) {
	body := struct {
		Filename string `json:"filename"`
		TmpKey   string `json:"tmpKey"`
		Uid      int    `json:"uid"`
		Class    string `json:"class"`
	}{}
	//接收回调参数

	if body.TmpKey == "" || len(body.TmpKey) != 16 || body.Class == "" || body.Uid == 0 {
		fmt.Println("参数错误")
		return
	}
	md5Str := utils.Md5Encode(fmt.Sprintf("%s|%d|%s|imguploadsafe", body.TmpKey[0:13], body.Uid, body.Class))
	if strings.ToUpper(md5Str[0:3]) != body.TmpKey[13:16] {
		fmt.Println("文件异常")
		return
	}
	ossUrl := "oss bucket url"

	imgSrc := ossUrl + body.Filename
	res := map[string]interface{}{
		"src": imgSrc,
	}
	fmt.Println(res)

	return
}
