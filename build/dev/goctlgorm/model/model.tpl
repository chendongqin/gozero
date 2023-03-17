package {{.pkg}}

import (
	"context"
	"gozore-mall/service/gormpool"
)

type Gorm{{.upperStartCamelObject}}Model struct {
	gormpool.GormBaseModel
}

func NewGorm{{.upperStartCamelObject}}Model(ctx context.Context, gormQuery *gormpool.CommonQuery) *Gorm{{.upperStartCamelObject}}Model {
	return &Gorm{{.upperStartCamelObject}}Model{
		GormBaseModel: gormpool.GormBaseModel{
			Ctx:       ctx,
			GormQuery: gormQuery,
		},
	}
}