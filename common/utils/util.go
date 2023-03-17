package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func IsChinesePhone(phone string) bool {
	regx := regexp.MustCompile("(^0086\\d+$)|(^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\\d{8}$)")
	return regx.MatchString(phone)
}

//GetRandom 生成随机字符串
func GetRandom(n int, strType int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$%^&*()+[]{}/<>;:=.,?"
	switch strType {
	case 1:
		letterBytes = "1234567890"
	case 2:
		letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case 3:
		letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	default:
	}
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

//RandomUid 生成随机字符串
func RandomUid(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func Md5Encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//InArrayString 判断字符串是否在字符串数组内
func InArrayString(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

//InArrayInt 判断字int是否在int数组内
func InArrayInt(need int, needArr []int) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

//InArrayInt64 判断字int64是否在int64数组内
func InArrayInt64(need int64, needArr []int64) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func ToInt(val interface{}) int {
	switch val.(type) {
	case float64:
		return int(val.(float64))
	case int32:
		return int(val.(int32))
	case int:
		return val.(int)
	case int64:
		return int(val.(int64))
	case uint8:
		return int(val.(uint8))
	case string:
		return ParseInt(val.(string), 0)
	}
	return 0
}

func ParseInt(b string, defInt int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return defInt
	} else {
		return id
	}
}

func ToInt64(val interface{}) int64 {
	switch val.(type) {
	case float64:
		return int64(val.(float64))
	case int:
		return int64(val.(int))
	case int32:
		return int64(val.(int32))
	case int64:
		return val.(int64)
	case string:
		return ParseInt64String(val.(string))
	}
	return 0
}

func ParseInt64String(b string) int64 {
	id, _ := strconv.ParseInt(b, 10, 64)
	return id
}

func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		return fmt.Sprint(result)
	}
}

func ToFloat64(value interface{}) float64 {
	num, err := strconv.ParseFloat(ToString(value), 64)
	if err != nil {
		return 0
	}
	return num
}

func ToFloat32(value interface{}) float32 {
	return float32(ToFloat64(value))
}

func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, e = net.DialTimeout(network, addr, time.Second*30) //设置建立连接超时
				if e != nil {
					return nil, e
				}
				err := conn.SetDeadline(time.Now().Add(time.Second * 30)) //设置发送接受数据超时
				return conn, err
			},
			ResponseHeaderTimeout: time.Second * 10,
		},
	}
}

func InterfaceToStruct(origin, scr interface{}) {
	byteData, _ := json.Marshal(origin)
	_ = json.Unmarshal(byteData, scr)
	return
}

//ZTOSign 中通签名
func ZTOSign(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	signData := h.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(signData)
	return sign
}

//ZTOClient  ztoUrl请求接口地址，param JSON字符串，ztoAppKey 中通AppKey， ztoAppSecret 中通秘钥， method 请求方式GET,POST
func ZTOClient(ztoUrl string, param string, ztoAppKey string, ztoAppSecret string, method string) (resp string, e error) {
	jsonParam := bytes.NewBuffer([]byte(param))
	//参数拼接秘钥 MD5后 base64
	sign := ZTOSign(param + ztoAppSecret)
	//创建Curl客户端
	client := NewHttpClient()
	if method == "GET" {
		//处理中
	} else {
		request, _ := http.NewRequest(method, ztoUrl, jsonParam)
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		request.Header.Set("x-appKey", ztoAppKey)
		request.Header.Set("x-dataDigest", sign)
		response, _ := client.Do(request)
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		resp = string(body)
		e = err
	}
	return resp, e
}

//KeyTransformation 入参的KEY转为mysql的key
func KeyTransformation(s string) (respString string) {
	for _, st := range s {
		if unicode.IsUpper(st) {
			respString = respString + "_" + strings.ToLower(string(st))
		}
		respString = respString + string(st)
	}
	return respString
}

func SortMap(params map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	newParams := map[string]interface{}{}
	for _, v := range keys {
		newParams[v] = params[v]
	}
	return newParams
}

func PriceFloat64(fv float64) float64 {
	val, err := strconv.ParseFloat(fmt.Sprintf("%.2f", fv), 64)
	if err != nil {
		return fv
	}
	return val
}

func RoundFloat64ToInt64(fv float64) int64 {
	vs := strings.Split(ToString(fv), ".")
	if len(vs) == 0 {
		return 0
	}
	val := ToInt64(vs[0])
	return val
}

func PriceFloat32(fv float32) float32 {
	val, err := strconv.ParseFloat(fmt.Sprintf("%.2f", fv), 64)
	if err != nil {
		return fv
	}
	return float32(val)
}

func Int64SliceJoin(arr []int64, sep string) string {
	newArr := make([]string, 0)
	for _, v := range arr {
		newArr = append(newArr, ToString(v))
	}
	return strings.Join(newArr, sep)
}

func IntSliceJoin(arr []int, sep string) string {
	newArr := make([]string, 0)
	for _, v := range arr {
		newArr = append(newArr, ToString(v))
	}
	return strings.Join(newArr, sep)
}
