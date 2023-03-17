package dict

import (
	"context"
	"gozore-mall/app/model/sysadmin"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/app/model"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddConfigDictLogic {
	return &AddConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddConfigDictLogic) AddConfigDict(req *types.AddConfigDictReq) error {
	var dictionary = &sysadmin.SysDictionary{}
	if req.ParentId != globalkey.SysTopParentId {
		err := l.svcCtx.GormQuery.FindOne(l.ctx, dictionary, req.ParentId)
		if err != nil {
			return errorx.NewDefaultError(errorx.ParentDictionaryIdErrorCode)
		}
	}
	err := l.svcCtx.GormQuery.FindOneBy(l.ctx, dictionary, "unique_key", req.UniqueKey)
	if err == model.ErrNotFound {
		err = copier.Copy(dictionary, req)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		err = l.svcCtx.GormQuery.Save(l.ctx, dictionary)
		if err != nil {
			return errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}

		return nil
	} else {

		return errorx.NewDefaultError(errorx.AddDictionaryErrorCode)
	}
}
