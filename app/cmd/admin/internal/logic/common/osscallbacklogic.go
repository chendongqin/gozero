package common

import (
	"context"
	"fmt"
	"gozore-mall/common/errorx"
	"gozore-mall/common/utils"
	"strings"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOssCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssCallbackLogic {
	return &OssCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssCallbackLogic) OssCallback(req *types.OssCallackReq) (resp *types.OssCallackResp, err error) {
	if req.TmpKey == "" || len(req.TmpKey) != 16 || req.Class == "" || req.UserId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	md5Str := utils.Md5Encode(fmt.Sprintf("%s|%d|%s|imguploadsafe", req.TmpKey[0:13], req.UserId, req.Class))
	if strings.ToUpper(md5Str[0:3]) != req.TmpKey[13:16] {
		return nil, errorx.NewDefaultError(errorx.OssCallbackErrorCode)
	}

	imgSrc := l.svcCtx.Config.Oss.OssUrl + req.Filename
	resp = &types.OssCallackResp{Src: imgSrc}

	return
}
