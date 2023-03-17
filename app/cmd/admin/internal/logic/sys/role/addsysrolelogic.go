package role

import (
	"context"
	"encoding/json"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/service/gormpool"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysRoleLogic {
	return &AddSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysRoleLogic) AddSysRole(req *types.AddSysRoleReq) error {
	var sysRole = new(sysadmin.SysRole)
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, sysRole, "unique_key", req.UniqueKey)
	if err == model.ErrNotFound {
		if req.ParentId != globalkey.SysTopParentId {
			var sysParentRole = new(sysadmin.SysRole)
			err = l.svcCtx.GormQuery.FindOne(l.ctx, sysParentRole, req.ParentId)
			if err != nil {
				return errorx.NewDefaultError(errorx.ParentRoleIdErrorCode)
			}
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
	} else {

		return errorx.NewDefaultError(errorx.AddRoleErrorCode)
	}
}
