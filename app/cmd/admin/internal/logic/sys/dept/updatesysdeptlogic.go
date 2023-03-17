package dept

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysDeptLogic {
	return &UpdateSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysDeptLogic) UpdateSysDept(req *types.UpdateSysDeptReq) error {
	dept := &sysadmin.SysDept{}
	if req.ParentId != globalkey.SysTopParentId {
		sysParentDept := &sysadmin.SysDept{}
		err := l.svcCtx.GormQuery.FindOne(l.ctx, sysParentDept, req.ParentId)
		if err != nil {
			return errorx.NewDefaultError(errorx.ParentDeptIdErrorCode)
		}
	}

	if req.Id == req.ParentId {
		return errorx.NewDefaultError(errorx.ParentDeptErrorCode)
	}

	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, dept, "unique_key", req.UniqueKey)
	if err != model.ErrNotFound && dept.Id != req.Id {
		return errorx.NewDefaultError(errorx.UpdateDeptUniqueKeyErrorCode)
	}

	deptIds := make([]int64, 0)
	deptIds = l.getSubDept(deptIds, req.Id)
	if utils.ArrayContainValue(deptIds, req.ParentId) {
		return errorx.NewDefaultError(errorx.SetParentIdErrorCode)
	}
	sysDept := &sysadmin.SysDept{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, sysDept, req.Id)
	if err != nil {
		return errorx.NewDefaultError(errorx.DeptIdErrorCode)
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
}

func (l *UpdateSysDeptLogic) getSubDept(deptIds []int64, id int64) []int64 {
	deptList := make([]sysadmin.SysDept, 0)
	err := l.svcCtx.GormQuery.QueryFindAll(l.ctx, &deptList, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{id},
	})
	if err != nil && err != model.ErrNotFound {
		return deptIds
	}

	for _, v := range deptList {
		deptIds = append(deptIds, v.Id)
		deptIds = l.getSubDept(deptIds, v.Id)
	}

	return deptIds
}
