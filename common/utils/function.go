package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"gozore-mall/common/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
)

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func GetAdminUserId(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value(globalkey.SysJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}

	return uid
}

func GetUserRoleIds(ctx context.Context) string {
	tel := ctx.Value(globalkey.SysRoleIds)
	return ToString(tel)
}

func GetReqJson(ctx context.Context, req interface{}) error {
	if jsonStr, ok := ctx.Value(globalkey.SysReqJson).(json.Number); ok {
		err := json.Unmarshal([]byte(jsonStr.String()), req)
		return err
	} else {
		return errors.New("参数错误")
	}
}

func SetReqJson(ctx context.Context, val interface{}) context.Context {
	reqJson, err := json.Marshal(val)
	fmt.Println(err)
	ctx = context.WithValue(ctx, globalkey.SysReqJson, reqJson)
	return ctx
}

func GetUserId(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value(globalkey.ApiJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}

	return uid
}

func ArrayUniqueValue[T any](arr []T) []T {
	size := len(arr)
	result := make([]T, 0, size)
	temp := map[any]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; ok != true {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}

	return result
}

func ArrayContainValue(arr []int64, search int64) bool {
	for _, v := range arr {
		if v == search {
			return true
		}
	}

	return false
}

func Intersect(slice1 []int64, slice2 []int64) []int64 {
	m := make(map[int64]int64)
	n := make([]int64, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}

	return n
}

func Difference(slice1 []int64, slice2 []int64) []int64 {
	m := make(map[int64]int)
	n := make([]int64, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, v := range slice1 {
		times, _ := m[v]
		if times == 0 {
			n = append(n, v)
		}
	}

	return n
}

func GetRemoteClientIp(r *http.Request) string {
	remoteIp := r.RemoteAddr

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteIp = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(remoteIp)
	}
	//本地ip
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}
	return remoteIp
}

func GetSource(r *http.Request) string {
	source := r.Header.Get("SOURCE")
	return source
}

func GetChannel(r *http.Request) string {
	channel := r.Header.Get("CHANNEL")
	return channel
}

func GetShareCode(r *http.Request) string {
	channel := r.Header.Get("SHARE_CODE")
	return channel
}

type RequestDeviceInfo struct {
	Device       string `json:"device"`
	IEMI         string `json:"iemi"`
	OAID         string `json:"oaid"`
	Ua           string `json:"ua"`
	Version      string `json:"version"`
	PhoneBrand   string `json:"phone_brand"`
	PhoneFrom    string `json:"phone_from"`
	PhoneBoard   string `json:"phone_board"`
	PhoneSn      string `json:"phone_sn"`
	PhoneVersion string `json:"phone_version"`
}

func GetDeviceInfo(r *http.Request) *RequestDeviceInfo {
	return &RequestDeviceInfo{
		Device:       r.Header.Get("DEVICE"),
		IEMI:         r.Header.Get("IEMI"),
		OAID:         r.Header.Get("OAID"),
		Ua:           r.Header.Get("UA"),
		Version:      r.Header.Get("VERSION"),
		PhoneBrand:   r.Header.Get("PHONE_BRAND"),
		PhoneFrom:    r.Header.Get("PHONE_FROM"),
		PhoneBoard:   r.Header.Get("PHONE_BOARD"),
		PhoneSn:      r.Header.Get("PHONE_SN"),
		PhoneVersion: r.Header.Get("PHONE_VERSION"),
	}
}

func ExportCsv(w *http.Response, downFilename string, header []string, data [][]string) ([]byte, error) {
	//内容先写入buffer缓存
	buff := new(bytes.Buffer)
	//写入UTF-8 BOM,此处如果不写入就会导致写入的汉字乱码
	buff.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(buff)
	wStr.Write(header)
	for _, s := range data {
		wStr.Write(s)
	}
	wStr.Flush()
	//指定下载文件名，可以注释掉，让前端处理文件名
	w.Header.Set("Content-Disposition", "attachment; filename="+downFilename)
	w.Header.Set("Content-Type", "text/csv") //设置为 .csv 格式文件
	return buff.Bytes(), nil
}
