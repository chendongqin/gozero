package menu

import (
	"context"
	"encoding/json"
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

type UpdateSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysPermMenuLogic {
	return &UpdateSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysPermMenuLogic) UpdateSysPermMenu(req *types.UpdateSysPermMenuReq) error {
	userId := utils.GetAdminUserId(l.ctx)
	if userId != globalkey.SysSuperUserId {
		gormSysUser := sysadmin.NewGormSysUserModel(l.ctx, l.svcCtx.GormQuery)
		if !gormSysUser.CheckSuperRole(userId) {
			return errorx.NewDefaultError(errorx.NotPermMenuErrorCode)
		}
	}

	if req.ParentId != globalkey.SysTopParentId {
		parentPermMenu := &sysadmin.SysPermMenu{}
		err := l.svcCtx.GormQuery.FindOne(l.ctx, parentPermMenu, req.ParentId)
		if err != nil {
			return errorx.NewDefaultError(errorx.ParentPermMenuIdErrorCode)
		}

		if parentPermMenu.Type == 2 {
			return errorx.NewDefaultError(errorx.SetParentTypeErrorCode)
		}
	}

	if req.Id <= globalkey.SysProtectPermMenuMaxId {
		return errorx.NewDefaultError(errorx.ForbiddenErrorCode)
	}

	if req.Id == req.ParentId {
		return errorx.NewDefaultError(errorx.ParentPermMenuErrorCode)
	}

	permMenuIds := make([]int64, 0)
	permMenuIds = l.getSubPermMenu(permMenuIds, req.Id)
	if utils.ArrayContainValue(permMenuIds, req.ParentId) {
		return errorx.NewDefaultError(errorx.SetParentIdErrorCode)
	}

	permMenu := &sysadmin.SysPermMenu{}
	err := l.svcCtx.GormQuery.FindOne(l.ctx, permMenu, req.Id)
	if err != nil {
		return errorx.NewDefaultError(errorx.PermMenuIdErrorCode)
	}

	err = copier.Copy(permMenu, req)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(req.Perms)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	permMenu.Perms = string(bytes)
	err = l.svcCtx.GormQuery.Save(l.ctx, permMenu)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}

func (l *UpdateSysPermMenuLogic) getSubPermMenu(permMenuIds []int64, id int64) []int64 {
	permMenuList := make([]sysadmin.SysPermMenu, 0)
	err := l.svcCtx.GormQuery.QueryFindAll(l.ctx, &permMenuList, gormpool.Conditions{
		Where:  "parent_id=?",
		Values: []interface{}{id},
	})
	if err != nil && err != model.ErrNotFound {
		return permMenuIds
	}

	for _, v := range permMenuList {
		permMenuIds = append(permMenuIds, v.Id)
		permMenuIds = l.getSubPermMenu(permMenuIds, v.Id)
	}

	return permMenuIds
}
