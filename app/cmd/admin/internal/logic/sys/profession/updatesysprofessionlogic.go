package profession

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

type UpdateSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysProfessionLogic {
	return &UpdateSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysProfessionLogic) UpdateSysProfession(req *types.UpdateSysProfessionReq) error {
	sysProfession := &sysadmin.SysProfession{}
	err := l.svcCtx.GormQuery.FindOne(l.ctx, sysProfession, req.Id)
	if err != nil {
		return errorx.NewDefaultError(errorx.ProfessionIdErrorCode)
	}

	if req.Status == globalkey.SysDisable {
		count, _ := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysUser{}, gormpool.Conditions{
			Where:  "profession_id=?",
			Values: []interface{}{req.Id},
		})
		if count > 0 {
			return errorx.NewDefaultError(errorx.JobIsUsingErrorCode)
		}
	}

	err = copier.Copy(sysProfession, req)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.GormQuery.Save(l.ctx, sysProfession)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
