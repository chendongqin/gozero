package user

import (
	"context"
	"strconv"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() error {
	userId := strconv.FormatInt(utils.GetAdminUserId(l.ctx), 10)
	//_, _ = l.svcCtx.Redis.Del(globalkey.SysAdminPermMenuCachePrefix + userId)
	_, _ = l.svcCtx.Redis.Del(globalkey.SysAdminOnlineUserCachePrefix + userId)
	//_, _ = l.svcCtx.Redis.Del(globalkey.SysUserIdCachePrefix + userId)

	return nil
}
