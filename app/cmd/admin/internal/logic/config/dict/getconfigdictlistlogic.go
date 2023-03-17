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

type GetConfigDictListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigDictListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigDictListLogic {
	return &GetConfigDictListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigDictListLogic) GetConfigDictList() (resp *types.ConfigDictListResp, err error) {
	configDictionaryList := make([]sysadmin.SysDictionary, 0)
	err = l.svcCtx.GormQuery.
		QueryFindAll(l.ctx, &configDictionaryList, gormpool.Conditions{
			Where:   "parent_id=0",
			OrderBy: "order_num DESC",
		})
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var dictionary types.ConfigDict
	dictionaryList := make([]types.ConfigDict, 0)
	for _, v := range configDictionaryList {
		err := copier.Copy(&dictionary, &v)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		dictionaryList = append(dictionaryList, dictionary)
	}

	return &types.ConfigDictListResp{
		List: dictionaryList,
	}, nil
}
