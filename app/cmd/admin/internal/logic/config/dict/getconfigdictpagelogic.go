package dict

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigDictPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigDictPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigDictPageLogic {
	return &GetConfigDictPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigDictPageLogic) GetConfigDictPage(req *types.ConfigDictPageReq) (resp *types.ConfigDictPageResp, err error) {
	configDictionaryList := make([]sysadmin.SysDictionary, 0)
	total, err := l.svcCtx.GormQuery.QueryPageListAndTotal(l.ctx, &sysadmin.SysDictionary{}, &configDictionaryList, gormpool.PageList{
		Page:   req.Page,
		Limit:  req.Limit,
		Where:  "parent_id=?",
		Values: []interface{}{req.ParentId},
	})
	var dictionary types.ConfigDict
	dictionaryList := make([]types.ConfigDict, 0)
	for _, sysDictionary := range configDictionaryList {
		err := copier.Copy(&dictionary, &sysDictionary)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		dictionaryList = append(dictionaryList, dictionary)
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.ConfigDictPageResp{
		List:       dictionaryList,
		Pagination: pagination,
	}, nil
}
