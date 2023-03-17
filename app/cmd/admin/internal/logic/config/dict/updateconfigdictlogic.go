package dict

import (
	"context"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConfigDictLogic {
	return &UpdateConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConfigDictLogic) UpdateConfigDict(req *types.UpdateConfigDictReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		var configDictionary = &sysadmin.SysDictionary{}
		err := l.svcCtx.GormQuery.FindOneBy(l.ctx, configDictionary, "parent_id", req.ParentId)
		if err != nil {
			return errorx.NewDefaultError(errorx.ParentDictionaryIdErrorCode)
		}
	}

	var configDictionary = &sysadmin.SysDictionary{}
	err := l.svcCtx.GormQuery.FindOne(l.ctx, configDictionary, req.Id)
	if err != nil {
		return errorx.NewDefaultError(errorx.DictionaryIdErrorCode)
	}

	err = copier.Copy(configDictionary, req)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.GormQuery.Save(l.ctx, configDictionary)
	if err != nil {
		return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	return nil
}
