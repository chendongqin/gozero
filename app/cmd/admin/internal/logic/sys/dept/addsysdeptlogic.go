package dept

import (
	"context"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysDeptLogic {
	return &AddSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysDeptLogic) AddSysDept(req *types.AddSysDeptReq) error {
	var sysDept = &sysadmin.SysDept{}
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, sysDept, "unique_key", req.UniqueKey)
	if err == model.ErrNotFound {
		if req.ParentId != globalkey.SysTopParentId {
			var sysParentDept = &sysadmin.SysDept{}
			err := l.svcCtx.GormQuery.FindOne(l.ctx, sysParentDept, req.ParentId)
			if err != nil {
				return errorx.NewDefaultError(errorx.ParentDeptIdErrorCode)
			}
		}

		err = copier.Copy(sysDept, req)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		err = l.svcCtx.GormQuery.Save(l.ctx, sysDept)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		return nil
	} else {
		return errorx.NewDefaultError(errorx.AddDeptErrorCode)
	}
}
