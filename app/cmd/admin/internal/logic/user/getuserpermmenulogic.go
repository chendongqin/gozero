package user

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/utils"
	"gozore-mall/service/gormpool"
	"strconv"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermMenuLogic {
	return &GetUserPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPermMenuLogic) GetUserPermMenu() (resp *types.UserPermMenuResp, err error) {
	userId := utils.GetAdminUserId(l.ctx)

	online, err := l.svcCtx.Redis.Get(globalkey.SysAdminOnlineUserCachePrefix + strconv.FormatInt(userId, 10))
	if err != nil || online == "" {
		return nil, errorx.NewDefaultError(errorx.AuthErrorCode)
	}

	// 查询用户信息
	user := &sysadmin.SysUser{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, user, userId)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var roles []int64
	// 用户所属角色
	err = json.Unmarshal([]byte(user.RoleIds), &roles)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var permMenu []int64
	var userPermMenu []sysadmin.SysPermMenu

	userPermMenu, permMenu, err = l.countUserPermMenu(roles, permMenu)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var menu types.Menu
	menuList := make([]types.Menu, 0)
	permList := make([]string, 0)
	//_, err = l.svcCtx.Redis.Del(globalkey.SysAdminPermMenuCachePrefix + strconv.FormatInt(userId, 10))
	for _, v := range userPermMenu {
		err := copier.Copy(&menu, &v)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		if menu.Type != globalkey.SysDefaultPermType {
			menuList = append(menuList, menu)
		}
		var permArray []string
		err = json.Unmarshal([]byte(v.Perms), &permArray)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		for _, p := range permArray {
			//_, err := l.svcCtx.Redis.Sadd(globalkey.SysAdminPermMenuCachePrefix+strconv.FormatInt(userId, 10), globalkey.SysPermMenuPrefix+p)
			//if err != nil {
			//	return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
			//}
			permList = append(permList, "/"+p)
		}

	}

	return &types.UserPermMenuResp{Menus: menuList, Perms: utils.ArrayUniqueValue[string](permList)}, nil
}

func (l *GetUserPermMenuLogic) countUserPermMenu(roles []int64, permMenu []int64) ([]sysadmin.SysPermMenu, []int64, error) {
	if utils.ArrayContainValue(roles, globalkey.SysSuperRoleId) {
		sysPermMenus := make([]sysadmin.SysPermMenu, 0)
		err := l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysPermMenus, gormpool.Conditions{
			OrderBy: "order_num DESC",
		})
		if err != nil {
			return nil, permMenu, err
		}

		return sysPermMenus, permMenu, nil
	} else {
		for _, roleId := range roles {
			// 查询角色信息
			role := &sysadmin.SysRole{}
			err := l.svcCtx.GormQuery.FindOne(l.ctx, role, roleId)
			if err != nil && err != model.ErrNotFound {
				return nil, permMenu, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
			}

			var perms []int64
			// 角色所拥有的权限id
			err = json.Unmarshal([]byte(role.PermMenuIds), &perms)
			if err != nil {
				return nil, permMenu, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
			}

			if role.Status != 0 {
				permMenu = append(permMenu, perms...)
			}
			// 汇总用户所属角色权限id
			permMenu = l.getRolePermMenu(permMenu, roleId)
		}

		// 过滤重复的权限id
		permMenu = utils.ArrayUniqueValue[int64](permMenu)
		var ids []int64
		for _, id := range permMenu {
			ids = append(ids, id)
		}

		if len(ids) == 0 {
			return nil, permMenu, nil
		}

		// 根据权限id获取具体权限
		sysPermMenus := make([]sysadmin.SysPermMenu, 0)
		err := l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysPermMenus, gormpool.Conditions{
			Where:  "id IN(?)",
			Values: []interface{}{ids},
		})
		if err != nil {
			return nil, permMenu, err
		}

		return sysPermMenus, permMenu, nil
	}
}

func (l *GetUserPermMenuLogic) getRolePermMenu(perms []int64, roleId int64) []int64 {
	roles := make([]sysadmin.SysRole, 0)
	err := l.svcCtx.GormQuery.QueryFindAll(l.ctx, &roles, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{roleId},
	})
	if err != nil && err != model.ErrNotFound {
		return perms
	}

	for _, role := range roles {
		var subPerms []int64
		err = json.Unmarshal([]byte(role.PermMenuIds), &subPerms)
		perms = append(perms, subPerms...)
		perms = l.getRolePermMenu(perms, role.Id)
	}

	return perms
}
