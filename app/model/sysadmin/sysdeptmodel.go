package sysadmin

import (
	"context"
	"gozore-mall/service/gormpool"
)

type GormSysDeptModel struct {
	gormpool.GormBaseModel
}

func NewGormSysDeptModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysDeptModel {
	return &GormSysDeptModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}
