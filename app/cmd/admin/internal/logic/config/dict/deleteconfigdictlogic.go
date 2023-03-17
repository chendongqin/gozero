package dict

import (
	"context"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConfigDictLogic {
	return &DeleteConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteConfigDictLogic) DeleteConfigDict(req *types.DeleteConfigDictReq) error {
	if req.Id <= globalkey.SysProtectDictionaryMaxId {
		return errorx.NewDefaultError(errorx.ForbiddenErrorCode)
	}

	var dictionary = &sysadmin.SysDictionary{}
	err := l.svcCtx.GormQuery.FindOne(l.ctx, dictionary, req.Id)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	if dictionary.Id > 0 {
		return errorx.NewDefaultError(errorx.DeleteDictionaryErrorCode)
	}

	err = l.svcCtx.GormQuery.DeleteByPk(l.ctx, &sysadmin.SysDictionary{}, req.Id)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
