package gormpool

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type GormBaseModel struct {
	Ctx       context.Context
	GormQuery *CommonQuery
	Page      int64
	Limit     int64
	OrderBy   []string
	AndWhere  []string
	AndValues []interface{}
	OrWhere   []string
	OrValues  []interface{}
}

func (m *GormBaseModel) SetPager(page, limit int64) {
	m.Page = page
	m.Limit = limit
	return
}

func (m *GormBaseModel) SetOrderBy(orderBy string) {
	m.OrderBy = append(m.OrderBy, orderBy)
	return
}

func (m *GormBaseModel) GetOrderBy() string {
	return strings.Join(m.OrderBy, ",")
}

func (m *GormBaseModel) SetAndWhereValue(where string, value ...interface{}) {
	m.AndWhere = append(m.AndWhere, where)
	m.AndValues = append(m.AndValues, value...)
	return
}

func (m *GormBaseModel) SetOrWhereValue(where string, value ...interface{}) {
	m.OrWhere = append(m.OrWhere, where)
	m.OrValues = append(m.OrValues, value...)
	return
}

func (m *GormBaseModel) SetAndWhere(where string) {
	m.AndWhere = append(m.AndWhere, where)
	return
}

func (m *GormBaseModel) SetOrWhere(where string) {
	m.OrWhere = append(m.OrWhere, where)
	return
}

func (m *GormBaseModel) SetAndValue(value ...interface{}) {
	m.AndValues = append(m.AndValues, value...)
	return
}

func (m *GormBaseModel) SetOrValue(value ...interface{}) {
	m.OrValues = append(m.OrValues, value...)
	return
}

func (m *GormBaseModel) GetWhere() (whereSql string, whereValues []interface{}) {
	defer m.initAndWhere()
	defer m.initOrWhere()
	sqlArr := make([]string, 0)
	if len(m.AndWhere) > 0 {
		sqlArr = append(sqlArr, fmt.Sprintf("(%s)", strings.Join(m.AndWhere, " AND ")))
		whereValues = append(whereValues, m.AndValues...)
	}
	if len(m.OrWhere) > 0 {
		sqlArr = append(sqlArr, fmt.Sprintf("(%s)", strings.Join(m.OrWhere, " OR ")))
		whereValues = append(whereValues, m.OrValues...)
	}
	if len(sqlArr) > 0 {
		if len(m.OrWhere) > 0 && len(m.AndWhere) > 0 {
			whereSql = fmt.Sprintf("(%s)", strings.Join(sqlArr, " AND "))
		} else {
			whereSql = fmt.Sprintf("%s", strings.Join(sqlArr, " AND "))
		}
	}

	return
}

func (m *GormBaseModel) GetAndWhere() (whereSql string, whereValues []interface{}) {
	defer m.initAndWhere()
	if len(m.AndWhere) > 0 {
		whereSql = fmt.Sprintf("(%s)", strings.Join(m.AndWhere, " AND "))
		whereValues = m.AndValues
	}

	return
}

func (m *GormBaseModel) GetOrWhere() (whereSql string, whereValues []interface{}) {
	defer m.initOrWhere()
	if len(m.OrWhere) > 0 {
		whereSql = fmt.Sprintf("(%s)", strings.Join(m.OrWhere, " OR "))
		whereValues = m.OrValues
	}

	return
}

//Init 重置
func (m *GormBaseModel) initAndWhere() {
	m.AndWhere = make([]string, 0)
	m.AndValues = make([]interface{}, 0)

	return
}

//Init 重置
func (m *GormBaseModel) initOrWhere() {
	m.OrWhere = make([]string, 0)
	m.OrValues = make([]interface{}, 0)

	return
}

func (m *GormBaseModel) NewGormSession() *gorm.Session {
	return &gorm.Session{NewDB: true}
}
