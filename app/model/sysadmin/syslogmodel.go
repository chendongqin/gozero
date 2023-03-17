package sysadmin

import (
	"context"
	"gozore-mall/service/gormpool"
)

type GormSysLogModel struct {
	gormpool.GormBaseModel
}

func NewGormSysLogModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysLogModel {
	return &GormSysLogModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}
