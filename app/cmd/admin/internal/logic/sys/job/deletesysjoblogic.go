package job

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysJobLogic {
	return &DeleteSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysJobLogic) DeleteSysJob(req *types.DeleteSysJobReq) error {
	count, _ := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysUser{}, gormpool.Conditions{
		Where:  "job_id=?",
		Values: []interface{}{req.Id},
	})
	if count != 0 {
		return errorx.NewDefaultError(errorx.DeleteJobErrorCode)
	}

	err := l.svcCtx.GormQuery.DeleteByPk(l.ctx, &sysadmin.SysJob{}, req.Id)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
