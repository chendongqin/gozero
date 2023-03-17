package menu

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSysPermMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysPermMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysPermMenuListLogic {
	return &GetSysPermMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysPermMenuListLogic) GetSysPermMenuList() (resp *types.SysPermMenuListResp, err error) {
	permMenus := make([]sysadmin.SysPermMenu, 0)
	err = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &permMenus, gormpool.Conditions{
		Where:  "",
		Values: nil,
	})
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	currentUserId := utils.GetAdminUserId(l.ctx)
	var currentUserPermMenuIds []int64
	if currentUserId != globalkey.SysSuperUserId {
		currentUserPermMenuIds = l.getCurrentUserPermMenuIds(currentUserId)
	}

	var menu types.PermMenu
	PermMenuList := make([]types.PermMenu, 0)
	for _, v := range permMenus {
		err := copier.Copy(&menu, &v)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		var perms []string
		err = json.Unmarshal([]byte(v.Perms), &perms)
		menu.Perms = perms
		if currentUserId == globalkey.SysSuperUserId {
			menu.Has = 1
		} else {
			if utils.ArrayContainValue(currentUserPermMenuIds, v.Id) {
				menu.Has = 1
			} else {
				menu.Has = 0
			}
		}
		PermMenuList = append(PermMenuList, menu)
	}

	return &types.SysPermMenuListResp{List: PermMenuList}, nil
}

func (l *GetSysPermMenuListLogic) getCurrentUserPermMenuIds(currentUserId int64) (ids []int64) {
	var currentPermMenuIds []int64
	if currentUserId != globalkey.SysSuperUserId {
		var currentUserRoleIds []int64
		currentUser := &sysadmin.SysUser{}
		_ = l.svcCtx.GormQuery.FindOne(l.ctx, currentUser, currentUserId)
		_ = json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
		sysRoles := make([]sysadmin.SysRole, 0)
		_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysRoles, gormpool.Conditions{
			Where:  "id!=? AND status=1 AND id IN(?)",
			Values: []interface{}{globalkey.SysSuperRoleId, currentUserRoleIds},
		})
		var rolePermMenus []int64
		for _, v := range sysRoles {
			err := json.Unmarshal([]byte(v.PermMenuIds), &rolePermMenus)
			if err != nil {
				return nil
			}
			currentPermMenuIds = append(currentPermMenuIds, rolePermMenus...)
		}
	}

	return currentPermMenuIds
}
