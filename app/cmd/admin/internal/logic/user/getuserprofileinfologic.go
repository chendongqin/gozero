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

type GetUserProfileInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserProfileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileInfoLogic {
	return &GetUserProfileInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserProfileInfoLogic) GetUserProfileInfo() (resp *types.UserProfileInfoResp, err error) {
	userId := utils.GetAdminUserId(l.ctx)
	user := &sysadmin.SysUser{}
	err = l.svcCtx.GormQuery.FindOne(l.ctx, user, userId)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return &types.UserProfileInfoResp{
		Username: user.Username,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Remark:   user.Remark,
		Avatar:   user.Avatar,
	}, nil
}
