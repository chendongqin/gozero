package dept

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysDeptLogic {
	return &DeleteSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysDeptLogic) DeleteSysDept(req *types.DeleteSysDeptReq) error {
	total, err := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysDept{}, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{req.Id},
	})
	if total != 0 {
		return errorx.NewDefaultError(errorx.DeleteDeptErrorCode)
	}

	count, err := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysUser{}, gormpool.Conditions{
		Where:  "dept_id=?",
		Values: []interface{}{req.Id},
	})
	if count != 0 {
		return errorx.NewDefaultError(errorx.DeptHasUserErrorCode)
	}

	err = l.svcCtx.GormQuery.DeleteByPk(l.ctx, &sysadmin.SysDept{}, req.Id)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
