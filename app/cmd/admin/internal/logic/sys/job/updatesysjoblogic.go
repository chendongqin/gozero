package job

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysJobLogic {
	return &UpdateSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysJobLogic) UpdateSysJob(req *types.UpdateSysJobReq) error {
	sysJob := &sysadmin.SysJob{}
	err := l.svcCtx.GormQuery.FindOne(l.ctx, sysJob, req.Id)
	if err != nil {
		return errorx.NewDefaultError(errorx.JobIdErrorCode)
	}

	if req.Status == globalkey.SysDisable {
		count, _ := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysUser{}, gormpool.Conditions{
			Where:  "job_id=?",
			Values: []interface{}{req.Id},
		})
		if count > 0 {
			return errorx.NewDefaultError(errorx.JobIsUsingErrorCode)
		}
	}

	err = copier.Copy(sysJob, req)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.GormQuery.Save(l.ctx, sysJob)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
