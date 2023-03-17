package menu

import (
	"context"
	"encoding/json"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysPermMenuLogic {
	return &AddSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysPermMenuLogic) AddSysPermMenu(req *types.AddSysPermMenuReq) error {
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

	var permMenu = new(sysadmin.SysPermMenu)
	err := copier.Copy(permMenu, req)
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
