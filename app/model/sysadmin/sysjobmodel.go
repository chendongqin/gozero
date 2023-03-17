package sysadmin

import (
	"context"
	"gozore-mall/service/gormpool"
)

type GormSysJobModel struct {
	gormpool.GormBaseModel
}

func NewGormSysJobModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysJobModel {
	return &GormSysJobModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}
