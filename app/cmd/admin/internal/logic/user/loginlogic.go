package user

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"net/http"
	"strconv"
	"time"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"
	"gozore-mall/common/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq, r *http.Request) (resp *types.LoginResp, err error) {
	verifyCode, _ := l.svcCtx.Redis.Get(globalkey.SysAdminLoginCaptchaCachePrefix + req.CaptchaId)
	if verifyCode != req.VerifyCode {
		return nil, errorx.NewDefaultError(errorx.CaptchaErrorCode)
	}

	sysUser := &sysadmin.SysUser{}
	err = l.svcCtx.GormQuery.FindOneBy(l.ctx, sysUser, "account", req.Account)
	if err != nil {
		return nil, errorx.NewDefaultError(errorx.AccountErrorCode)
	}

	if sysUser.Password != utils.MD5(req.Password+sysUser.Salt+l.svcCtx.Config.Salt) {
		return nil, errorx.NewDefaultError(errorx.PasswordErrorCode)
	}

	if sysUser.Status != globalkey.SysEnable {
		return nil, errorx.NewDefaultError(errorx.AccountDisableErrorCode)
	}

	if sysUser.Id != globalkey.SysSuperUserId {
		dept := &sysadmin.SysDept{}
		_ = l.svcCtx.GormQuery.FindOne(l.ctx, dept, sysUser.DeptId)
		if dept.Status == globalkey.SysDisable {
			return nil, errorx.NewDefaultError(errorx.AccountDisableErrorCode)
		}
	}

	token, _ := l.getJwtToken(sysUser.Id, sysUser.RoleIds)
	_, err = l.svcCtx.Redis.Del(req.CaptchaId)

	loginLog := &sysadmin.SysLog{
		UserId: sysUser.Id,
		Ip:     utils.GetRemoteClientIp(r),
		Uri:    r.RequestURI,
		Type:   1,
		Status: 1,
	}

	err = l.svcCtx.GormQuery.Save(l.ctx, loginLog)

	err = l.svcCtx.Redis.Setex(globalkey.SysAdminOnlineUserCachePrefix+strconv.FormatInt(sysUser.Id, 10), token, int(l.svcCtx.Config.JwtAuth.AccessExpire))
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return &types.LoginResp{
		Token: token,
	}, nil
}

func (l *LoginLogic) getJwtToken(userId int64, roleIds string) (string, error) {
	iat := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + l.svcCtx.Config.JwtAuth.AccessExpire
	claims["iat"] = iat
	claims[globalkey.SysJwtUserId] = userId
	claims[globalkey.SysRoleIds] = roleIds
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(l.svcCtx.Config.JwtAuth.AccessSecret))
}
