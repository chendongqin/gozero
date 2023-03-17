package profession

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysProfessionLogic {
	return &DeleteSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysProfessionLogic) DeleteSysProfession(req *types.DeleteSysProfessionReq) error {
	count, _ := l.svcCtx.GormQuery.CountTotal(l.ctx, &sysadmin.SysUser{}, gormpool.Conditions{
		Where:  "profession_id=?",
		Values: []interface{}{req.Id},
	})
	if count != 0 {
		return errorx.NewDefaultError(errorx.DeleteProfessionErrorCode)
	}

	err := l.svcCtx.GormQuery.DeleteByPk(l.ctx, &sysadmin.SysProfession{}, req.Id)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
