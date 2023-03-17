package common

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/utils"
	"gozore-mall/service/ali/oss"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOssTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOssTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOssTokenLogic {
	return &GetOssTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOssTokenLogic) GetOssToken(req *types.GetOssTokenReq) (resp *types.GetOssTokenResp, err error) {
	userId := utils.GetAdminUserId(l.ctx)
	className := "admin"
	uploadUrl := fmt.Sprintf("%s/%s/%d-%s.%s", className, time.Now().Format("20060102"), userId, utils.GetRandom(10, 3), req.FileType)
	randKey := utils.GetRandom(13, 2)
	md5Str := utils.Md5Encode(fmt.Sprintf("%s|%d|%s|imguploadsafe", randKey, userId, className))
	randKey = randKey + strings.ToUpper(md5Str[0:3])
	ossSrv := oss.NewOssSvr(l.svcCtx.Config.Oss.CallbackUrl, l.svcCtx.Config.Oss.AccessKeyId, l.svcCtx.Config.Oss.AccessKeySecret, l.svcCtx.Config.Oss.Endpoint)
	policyToken := ossSrv.GetPolicyTokenById(uploadUrl, randKey, className, userId)

	var policyTokenResp types.PolicyToken
	err = copier.Copy(&policyTokenResp, &policyToken)
	if err != nil {
		return nil, errorx.NewDefaultError(errorx.ServerErrorCode)
	}

	resp = &types.GetOssTokenResp{
		PolicyToken: policyTokenResp,
		OssUrl:      l.svcCtx.Config.Oss.OssUrl,
		BucketName:  l.svcCtx.Config.Oss.Bucket,
		Endpoint:    l.svcCtx.Config.Oss.Endpoint,
	}

	return
}
