package user

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"
	"strconv"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysUserLogic {
	return &UpdateSysUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysUserLogic) UpdateSysUser(req *types.UpdateSysUserReq) error {
	currentUserId := utils.GetAdminUserId(l.ctx)
	var currentUserRoleIds []int64
	var roleIds []int64
	if currentUserId == globalkey.SysSuperUserId {
		sysRoleList := make([]sysadmin.SysRole, 0)
		_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysRoleList, gormpool.Conditions{
			OrderBy: "order_num DESC",
		})
		for _, role := range sysRoleList {
			currentUserRoleIds = append(currentUserRoleIds, role.Id)
			roleIds = append(roleIds, role.Id)
		}

	} else {
		currentUser := &sysadmin.SysUser{}
		_ = l.svcCtx.GormQuery.FindOne(l.ctx, currentUser, currentUserId)
		err := json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		roleIds = append(roleIds, currentUserRoleIds...)
	}

	editUser := &sysadmin.SysUser{}
	err := l.svcCtx.GormQuery.FindOne(l.ctx, editUser, req.Id)
	if err != nil {
		return errorx.NewDefaultError(errorx.UserIdErrorCode)
	}

	var editUserRoleIds []int64
	err = json.Unmarshal([]byte(editUser.RoleIds), &editUserRoleIds)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}
	roleIds = append(roleIds, editUserRoleIds...)

	for _, id := range req.RoleIds {
		if !utils.ArrayContainValue(roleIds, id) {
			return errorx.NewDefaultError(errorx.AssigningRolesErrorCode)
		}
	}

	for _, id := range utils.Difference(editUserRoleIds, currentUserRoleIds) {
		if !utils.ArrayContainValue(req.RoleIds, id) {
			return errorx.NewDefaultError(errorx.AssigningRolesErrorCode)
		}
	}

	sysDept := &sysadmin.SysDept{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, sysDept, req.DeptId)
	if err != nil {
		return errorx.NewDefaultError(errorx.DeptIdErrorCode)
	}
	sysProfession := &sysadmin.SysProfession{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, sysProfession, req.ProfessionId)
	if err != nil {
		return errorx.NewDefaultError(errorx.ProfessionIdErrorCode)
	}

	sysJob := &sysadmin.SysJob{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, sysJob, req.JobId)
	if err != nil {
		return errorx.NewDefaultError(errorx.JobIdErrorCode)
	}

	err = copier.Copy(editUser, req)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(req.RoleIds)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	_, err = l.svcCtx.Redis.Del(globalkey.SysAdminOnlineUserCachePrefix + strconv.FormatInt(editUser.Id, 10))
	editUser.RoleIds = string(bytes)
	err = l.svcCtx.GormQuery.Save(l.ctx, editUser)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
