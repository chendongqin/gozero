package role

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/utils"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysRoleLogic {
	return &DeleteSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysRoleLogic) DeleteSysRole(req *types.DeleteSysRoleReq) error {
	if req.Id == globalkey.SysSuperRoleId {
		return errorx.NewDefaultError(errorx.ForbiddenErrorCode)
	}

	total, _ := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysRole{}, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{req.Id},
	})
	if total > 0 {
		return errorx.NewDefaultError(errorx.DeleteRoleErrorCode)
	}

	count, _ := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysUser{}, gormpool.Conditions{
		Where:  "JSON_CONTAINS(role_ids,JSON_ARRAY(?))",
		Values: []interface{}{req.Id},
	})
	if count != 0 {
		return errorx.NewDefaultError(errorx.RoleIsUsingErrorCode)
	}

	err := l.svcCtx.GormQuery.DeleteByPk(l.ctx, &sysadmin.SysRole{}, req.Id)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}
	_, _ = l.svcCtx.Redis.Del(globalkey.SysAdminPermMenuCachePrefix + utils.ToString(req.Id))

	return nil
}
