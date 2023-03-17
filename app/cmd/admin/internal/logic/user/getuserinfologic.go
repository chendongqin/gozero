package user

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/common/utils"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	userId := utils.GetAdminUserId(l.ctx)
	user := &sysadmin.SysUser{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, user, userId)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return &types.UserInfoResp{
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
