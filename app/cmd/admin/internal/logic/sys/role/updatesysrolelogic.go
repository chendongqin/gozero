package role

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"
	"strconv"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysRoleLogic {
	return &UpdateSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysRoleLogic) UpdateSysRole(req *types.UpdateSysRoleReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		sysParentRole := &sysadmin.SysRole{}
		err := l.svcCtx.GormQuery.FindOne(l.ctx, sysParentRole, req.ParentId)
		if err != nil {
			return errorx.NewDefaultError(errorx.ParentRoleIdErrorCode)
		}
	}

	if req.Id == globalkey.SysSuperRoleId {
		return errorx.NewDefaultError(errorx.NotPermMenuErrorCode)
	}

	if req.Id == req.ParentId {
		return errorx.NewDefaultError(errorx.ParentRoleErrorCode)
	}
	role := &sysadmin.SysRole{}
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, role, "unique_key", req.UniqueKey)
	if err != model.ErrNotFound && role.Id != req.Id {
		return errorx.NewDefaultError(errorx.UpdateRoleUniqueKeyErrorCode)
	}

	roleIds := make([]int64, 0)
	roleIds = l.getSubRole(roleIds, req.Id)
	if utils.ArrayContainValue(roleIds, req.ParentId) {
		return errorx.NewDefaultError(errorx.SetParentIdErrorCode)
	}
	sysRole := &sysadmin.SysRole{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, sysRole, req.Id)
	if err != nil {
		return errorx.NewDefaultError(errorx.RoleIdErrorCode)
	}

	err = copier.Copy(sysRole, req)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(req.PermMenuIds)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	sysRole.PermMenuIds = string(bytes)
	err = l.svcCtx.GormQuery.Save(l.ctx, sysRole)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	_, _ = l.svcCtx.Redis.Del(globalkey.SysAdminPermMenuCachePrefix + utils.ToString(sysRole.Id))
	// 根据权限id获取具体权限
	sysPermMenus := make([]sysadmin.SysPermMenu, 0)
	_ = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysPermMenus, gormpool.Conditions{
		Where:  "id IN(?)",
		Values: []interface{}{req.PermMenuIds},
	})
	for _, v := range sysPermMenus {
		var permArray []string
		_ = json.Unmarshal([]byte(v.Perms), &permArray)
		for _, p := range permArray {
			_, err = l.svcCtx.Redis.Sadd(globalkey.SysAdminPermMenuCachePrefix+strconv.FormatInt(sysRole.Id, 10), globalkey.SysPermMenuPrefix+p)
			if err != nil {
				return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
			}
		}
	}

	return nil
}

func (l *UpdateSysRoleLogic) getSubRole(roleIds []int64, id int64) []int64 {
	roleList := make([]sysadmin.SysRole, 0)
	err := l.svcCtx.GormQuery.QueryFindAll(l.ctx, &roleList, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{id},
	})
	if err != nil && err != model.ErrNotFound {
		return roleIds
	}

	for _, v := range roleList {
		roleIds = append(roleIds, v.Id)
		roleIds = l.getSubRole(roleIds, v.Id)
	}

	return roleIds
}
