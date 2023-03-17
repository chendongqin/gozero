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

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysPermMenuLogic {
	return &DeleteSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysPermMenuLogic) DeleteSysPermMenu(req *types.DeleteSysPermMenuReq) error {
	userId := utils.GetAdminUserId(l.ctx)
	if userId != globalkey.SysSuperUserId {
		gormSysUser := sysadmin.NewGormSysUserModel(l.ctx, l.svcCtx.GormQuery)
		if !gormSysUser.CheckSuperRole(userId) {
			return errorx.NewDefaultError(errorx.NotPermMenuErrorCode)
		}

	}
	//currentUserId := utils.GetAdminUserId(l.ctx)
	//if currentUserId != globalkey.SysSuperUserId {
	//	var currentUserPermMenuIds []int64
	//	currentUserPermMenuIds = l.getCurrentUserPermMenuIds(currentUserId)
	//	if !utils.ArrayContainValue(currentUserPermMenuIds, req.Id) {
	//		return errorx.NewDefaultError(errorx.NotPermMenuErrorCode)
	//	}
	//}

	if req.Id <= globalkey.SysProtectPermMenuMaxId {
		return errorx.NewDefaultError(errorx.ForbiddenErrorCode)
	}

	count, _ := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysPermMenu{}, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{req.Id},
	})
	if count != 0 {
		return errorx.NewDefaultError(errorx.DeletePermMenuErrorCode)
	}

	err := l.svcCtx.GormQuery.DeleteByPk(l.ctx, &sysadmin.SysPermMenu{}, req.Id)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}

func (l *DeleteSysPermMenuLogic) getCurrentUserPermMenuIds(currentUserId int64) (ids []int64) {
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
