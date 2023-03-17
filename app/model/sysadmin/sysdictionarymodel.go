package sysadmin

import (
	"context"
	"gozore-mall/service/gormpool"
)

type GormSysDictionaryModel struct {
	gormpool.GormBaseModel
}

func NewGormSysDictionaryModel(ctx context.Context, gormQuery *gormpool.CommonQuery) *GormSysDictionaryModel {
	return &GormSysDictionaryModel{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}
