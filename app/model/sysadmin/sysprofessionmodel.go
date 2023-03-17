package sysadmin

import (
	"context"
	"gozore-mall/service/gormpool"
)

type GormSysProfessionModel struct {
	gormpool.GormBaseModel
}

func NewGormSysProfessionModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysProfessionModel {
	return &GormSysProfessionModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}
