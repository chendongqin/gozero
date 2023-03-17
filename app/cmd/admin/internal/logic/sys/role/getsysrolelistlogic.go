package role

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSysRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysRoleListLogic {
	return &GetSysRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysRoleListLogic) GetSysRoleList() (resp *types.SysRoleListResp, err error) {
	sysRoleList := make([]sysadmin.SysRole, 0)
	err = l.svcCtx.GormQuery.QueryFindAll(l.ctx, &sysRoleList, gormpool.Conditions{
		OrderBy: "order_num DESC",
	})
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var role types.Role
	roleList := make([]types.Role, 0)
	for _, v := range sysRoleList {
		err := copier.Copy(&role, &v)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		var permMenuIds []int64
		err = json.Unmarshal([]byte(v.PermMenuIds), &permMenuIds)
		role.PermMenuIds = permMenuIds
		roleList = append(roleList, role)
	}

	return &types.SysRoleListResp{
		List: roleList,
	}, nil
}
