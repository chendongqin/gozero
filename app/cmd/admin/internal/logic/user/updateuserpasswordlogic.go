package user

import (
	"context"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdatePasswordReq) error {
	dictionary := &sysadmin.SysDictionary{}
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, dictionary, "unique_key", "sys_ch_pwd")
	if dictionary.Status == globalkey.SysDisable {
		return errorx.NewDefaultError(errorx.ForbiddenErrorCode)
	}

	userId := utils.GetAdminUserId(l.ctx)
	user := &sysadmin.SysUser{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, user, userId)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	if user.Password != utils.MD5(req.OldPassword+user.Salt+l.svcCtx.Config.Salt) {
		return errorx.NewDefaultError(errorx.PasswordErrorCode)
	}

	user.Password = utils.MD5(req.NewPassword + user.Salt + l.svcCtx.Config.Salt)
	err = l.svcCtx.GormQuery.Save(l.ctx, user)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
