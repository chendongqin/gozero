package profession

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

type GetSysProfessionPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysProfessionPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysProfessionPageLogic {
	return &GetSysProfessionPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysProfessionPageLogic) GetSysProfessionPage(req *types.SysProfessionPageReq) (resp *types.SysProfessionPageResp, err error) {
	sysProfessionList := make([]sysadmin.SysProfession, 0)
	total, err := l.svcCtx.GormQuery.QueryPageListAndTotal(l.ctx, &sysadmin.SysProfession{}, &sysProfessionList, gormpool.PageList{
		Page:    req.Page,
		Limit:   req.Limit,
		OrderBy: "order_num DESC",
	})
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var profession types.Profession
	professionList := make([]types.Profession, 0)
	for _, v := range sysProfessionList {
		err = copier.Copy(&profession, &v)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		professionList = append(professionList, profession)
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.SysProfessionPageResp{
		List:       professionList,
		Pagination: pagination,
	}, nil
}
