package profession

import (
	"context"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysProfessionLogic {
	return &AddSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysProfessionLogic) AddSysProfession(req *types.AddSysProfessionReq) error {
	var sysProfession = new(sysadmin.SysProfession)
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, sysProfession, "name", req.Name)
	if err == model.ErrNotFound {

		err = copier.Copy(sysProfession, req)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		err = l.svcCtx.GormQuery.Save(l.ctx, sysProfession)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		return nil
	} else {

		return errorx.NewDefaultError(errorx.AddProfessionErrorCode)
	}
}
