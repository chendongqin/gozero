package sysadmin

import (
	"context"
	"gozore-mall/service/gormpool"
)

type GormSysPermMenuModel struct {
	gormpool.GormBaseModel
}

func NewGormSysPermMenuModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysPermMenuModel {
	return &GormSysPermMenuModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}
