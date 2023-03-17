package sysadmin

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"
	"gozore-mall/service/gormpool"
	"strconv"
)

type GormSysRoleModel struct {
	gormpool.GormBaseModel
}

func NewGormSysRoleModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysRoleModel {
	return &GormSysRoleModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}

//SyncRolePermMenu 同步角色菜单
func (g *GormSysRoleModel) SyncRolePermMenu(redisClient *redis.Redis, roleId int64) error {
	if roleId == 0 {
		return nil
	}
	role := &SysRole{}
	err := g.GormQuery.FindOne(g.Ctx, role, roleId)
	if err != model.ErrNotFound {
		return nil
	}
	permMenuIds := make([]int64, 0)
	_ = json.Unmarshal([]byte(role.PermMenuIds), &permMenuIds)
	if len(permMenuIds) == 0 {
		return nil
	}
	_, _ = redisClient.Del(globalkey.SysAdminPermMenuCachePrefix + utils.ToString(roleId))
	// 根据权限id获取具体权限
	sysPermMenus := make([]SysPermMenu, 0)
	_ = g.GormQuery.QueryFindAll(g.Ctx, &sysPermMenus, gormpool.Conditions{
		Where:  "id IN(?)",
		Values: []interface{}{permMenuIds},
	})
	for _, v := range sysPermMenus {
		var permArray []string
		_ = json.Unmarshal([]byte(v.Perms), &permArray)
		for _, p := range permArray {
			_, err = redisClient.Sadd(globalkey.SysAdminPermMenuCachePrefix+strconv.FormatInt(role.Id, 10), globalkey.SysPermMenuPrefix+p)
			if err != nil {
				return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
			}
		}
	}

	return nil
}
