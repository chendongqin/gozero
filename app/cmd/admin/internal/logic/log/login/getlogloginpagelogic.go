package login

import (
	"context"
	"gozore-mall/app/model/sysadmin"
	"gozore-mall/service/gormpool"

	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogLoginPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogLoginPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogLoginPageLogic {
	return &GetLogLoginPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogLoginPageLogic) GetLogLoginPage(req *types.LogLoginPageReq) (resp *types.LogLoginPageResp, err error) {
	loginLogList := make([]sysadmin.SysLog, 0)
	total, err := l.svcCtx.GormQuery.QueryPageListAndTotal(l.ctx, &sysadmin.SysLog{}, &loginLogList, gormpool.PageList{
		Page:   req.Page,
		Limit:  req.Limit,
		Where:  "type=?",
		Values: []interface{}{globalkey.SysLoginLogType},
	})
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
	}

	var loginLog types.LogLogin
	logList := make([]types.LogLogin, 0)
	for _, v := range loginLogList {
		err := copier.Copy(&loginLog, &v)
		loginLog.CreateTime = v.Created.Format(globalkey.SysDateFormat)
		if err != nil {
			return nil, errorx.NewSystemError(errorx.ServerErrorCode, err.Error())
		}
		logList = append(logList, loginLog)
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.LogLoginPageResp{
		List:       logList,
		Pagination: pagination,
	}, nil
}
