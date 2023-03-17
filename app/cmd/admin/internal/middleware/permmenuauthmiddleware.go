package middleware

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"
	"gozore-mall/service/gormpool"
	"net/http"
	"strconv"
	"strings"
)

type PermMenuAuthMiddleware struct {
	Redis       *redis.Redis
	CommonQuery *gormpool.CommonQuery
}

func NewPermMenuAuthMiddleware(r *redis.Redis, commonQuery *gormpool.CommonQuery) *PermMenuAuthMiddleware {
	return &PermMenuAuthMiddleware{
		Redis:       r,
		CommonQuery: commonQuery,
	}
}

func (m *PermMenuAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if len(authToken) > 0 {
			userId := utils.GetAdminUserId(r.Context())
			online, err := m.Redis.Get(globalkey.SysAdminOnlineUserCachePrefix + strconv.FormatInt(userId, 10))
			if err != nil || online != authToken {
				httpx.Error(w, errorx.NewDefaultError(errorx.AuthErrorCode))
			} else {
				roleIdsStr := utils.GetUserRoleIds(r.Context())
				uri := strings.Split(r.RequestURI, "?")
				rolesIds := make([]int64, 0)
				_ = json.Unmarshal([]byte(roleIdsStr), &rolesIds)
				check := false
				for _, roleId := range rolesIds {
					if globalkey.SysSuperRoleId == roleId {
						check = true
						break
					}
					roleKey := globalkey.SysAdminPermMenuCachePrefix + strconv.FormatInt(roleId, 10)
					res, err := m.Redis.Smembers(roleKey)
					if err == nil && len(res) == 0 {
						sysRoleGormModel := sysadmin.NewGormSysRoleModel(r.Context(), m.CommonQuery)
						_ = sysRoleGormModel.SyncRolePermMenu(m.Redis, roleId)
					}
					is, err := m.Redis.Sismember(roleKey, uri[0])
					if err == nil && is == true {
						check = true
					}
				}
				if !check {
					httpx.Error(w, errorx.NewDefaultError(errorx.NotPermMenuErrorCode))
				} else {
					next(w, r)
				}
			}
		}
	}
}
