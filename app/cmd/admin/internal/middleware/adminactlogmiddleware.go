package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/utils"
	"gozore-mall/service/gormpool"
	"io/ioutil"
	"net/http"
	"strings"
)

var LogEncryptionWords = []string{"password"}

type AdminActLogMiddleware struct {
	CommonQuery *gormpool.CommonQuery
}

func NewAdminActLogMiddleware(commonQuery *gormpool.CommonQuery) *AdminActLogMiddleware {
	return &AdminActLogMiddleware{
		CommonQuery: commonQuery,
	}
}

func (m *AdminActLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//next(w, r)
		if len(r.Header.Get("Authorization")) > 0 {
			userId := utils.GetAdminUserId(r.Context())
			var body []byte
			if r.Body != nil {
				body, _ = ioutil.ReadAll(r.Body)
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			jsonBody := strings.Replace(strings.Replace(string(body), " ", "", -1), "\n", "", -1)
			bodyData := map[string]interface{}{}
			_ = json.Unmarshal([]byte(jsonBody), &bodyData)
			r = r.WithContext(utils.SetReqJson(r.Context(), bodyData))
			if userId > 0 {
				go func(gormQuery *gormpool.CommonQuery, userId int64, bodyData map[string]interface{}, r *http.Request) {
					for k, v := range bodyData {
						if utils.InArrayString(k, LogEncryptionWords) {
							switch v.(type) {
							case string:
								bodyData[k] = "******"
							}
						}
					}
					newBody, _ := json.Marshal(bodyData)
					loginLog := &sysadmin.SysLog{
						UserId:  userId,
						Ip:      utils.GetRemoteClientIp(r),
						Uri:     r.RequestURI,
						Request: string(newBody),
						Type:    2,
						Status:  1,
					}

					_ = gormQuery.Save(context.Background(), loginLog)
				}(m.CommonQuery, userId, bodyData, r)
			}
		}
		next(w, r)
	}
}
