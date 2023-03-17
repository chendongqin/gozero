package user

import (
	"context"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(req *types.UpdateProfileReq) error {
	dictionary := &sysadmin.SysDictionary{}
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, dictionary, "unique_key", "sys_userinfo")
	if dictionary.Status == globalkey.SysDisable {
		return errorx.NewDefaultError(errorx.ForbiddenErrorCode)
	}

	userId := utils.GetAdminUserId(l.ctx)
	user := &sysadmin.SysUser{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, user, userId)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	err = copier.Copy(user, req)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.GormQuery.Save(l.ctx, user)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
