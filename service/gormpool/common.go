package gormpool

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"gozore-mall/common/commonconst"
	"gozore-mall/common/utils"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type CommonTabler interface {
	TableName() string
	GetPKColumnName() string
	GetPkValue() interface{}
	SetTime()
}

type CacheSignalChan chan string

var (
	PK               = "id"
	CacheTableRowKey = "sql.cache.%s.%s.%v"
	CachePrefix      = "sql.md5."
)

// CacheSignal 信号检测
type CacheSignal struct {
	CacheOff     bool
	CacheErr     error
	CacheErrFile string
	CacheChan    CacheSignalChan
}

type CommonQuery struct {
	GormDb      *gorm.DB
	WriteGormDb *gorm.DB
	CacheRedis  cache.Cache
	CacheSignal *CacheSignal
	cacheQuery  bool
}

func NewCommonQuery(gormDb *gorm.DB, writeGormDb *gorm.DB, cacheRedis cache.Cache, cacheSignal *CacheSignal) *CommonQuery {

	return &CommonQuery{
		GormDb:      gormDb,
		WriteGormDb: writeGormDb,
		CacheRedis:  cacheRedis,
		CacheSignal: cacheSignal,
	}
}

func (m *CommonQuery) NewDbWithContext(ctx context.Context, gormSession *gorm.Session) *gorm.DB {
	return m.GormDb.Session(gormSession).WithContext(ctx)
}

func (m *CommonQuery) NewWriteDbWithContext(ctx context.Context, gormSession *gorm.Session) *gorm.DB {
	return m.WriteGormDb.Session(gormSession).WithContext(ctx)
}

func (m *CommonQuery) QueryPageListAndTotal(ctx context.Context, tableModel CommonTabler, list interface{}, pageList PageList) (int64, error) {
	var total = int64(0)
	db := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Model(tableModel).
		Where(pageList.Where, pageList.Values...)

	err := db.Count(&total).Error

	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 0, nil
	}

	start := (pageList.Page - 1) * pageList.Limit
	err = m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Select(pageList.Columns).
		Where(pageList.Where, pageList.Values...).
		Order(pageList.OrderBy).
		Offset(int(start)).
		Limit(int(pageList.Limit)).
		Find(list).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (m *CommonQuery) QueryList(ctx context.Context, list interface{}, pageList PageList) error {
	start := (pageList.Page - 1) * pageList.Limit
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Select(pageList.Columns).
		Where(pageList.Where, pageList.Values...).
		Order(pageList.OrderBy).
		Offset(int(start)).
		Limit(int(pageList.Limit)).
		Find(list).Error

	if err != nil {
		return err
	}

	return nil
}

func (m *CommonQuery) QueryFindAll(ctx context.Context, list interface{}, conditions Conditions) error {
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Select(conditions.Columns).
		Where(conditions.Where, conditions.Values...).
		Order(conditions.OrderBy).
		Find(list).Error

	return err
}

func (m *CommonQuery) QueryFindAllByTableName(ctx context.Context, tableName string, list interface{}, conditions Conditions) error {
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Table(tableName).
		Select(conditions.Columns).
		Where(conditions.Where, conditions.Values...).
		Order(conditions.OrderBy).
		Find(list).Error

	return err
}

func (m *CommonQuery) CountTotal(ctx context.Context, tableModel CommonTabler, conditions Conditions) (int64, error) {
	db := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Model(tableModel).
		Where(conditions.Where, conditions.Values...)
	var total = int64(0)
	err := db.Count(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (m *CommonQuery) Sum(ctx context.Context, tableModel CommonTabler, column string, conditions Conditions) (float64, error) {
	list := make([]struct {
		SumCount float64 `json:"sum_count"`
	}, 0)
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Model(tableModel).
		Select(fmt.Sprintf("SUM(`%s`) AS sum_count", column)).
		Where(conditions.Where, conditions.Values...).
		Find(list).Error

	if err != nil {
		return 0, err
	}

	if len(list) == 0 {
		return 0, err
	}

	return list[0].SumCount, nil
}

func (m *CommonQuery) CountTotalByTableName(ctx context.Context, tableName string, conditions Conditions) (int64, error) {
	db := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Table(tableName).
		Where(conditions.Where, conditions.Values...)
	var total = int64(0)
	err := db.Count(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (m *CommonQuery) FindOneWhere(ctx context.Context, tableModel CommonTabler, conditions Conditions) error {
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Select(conditions.Columns).
		Where(conditions.Where, conditions.Values...).
		Order(conditions.OrderBy).
		First(tableModel).Error

	return err
}

func (m *CommonQuery) FindSumWhere(ctx context.Context, tableModel CommonTabler, field string, totalValue *[]float64, conditions Conditions) error {
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Table(tableModel.TableName()).
		Where(conditions.Where, conditions.Values...).
		Pluck(field, &totalValue).
		//Scan(&totalValue).
		Error
	return err
}

//Create 新建
func (m *CommonQuery) Create(ctx context.Context, tableModel CommonTabler) error {
	defer m.checkAble()
	tableModel.SetTime()
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Create(tableModel).Error
	if err != nil {
		return err
	}

	return nil
}

//CreateAll 批量新建
func (m *CommonQuery) CreateAll(ctx context.Context, tableModels []CommonTabler) error {
	defer m.checkAble()
	if len(tableModels) == 0 {
		return nil
	}
	for k, tableModel := range tableModels {
		tableModel.SetTime()
		tableModels[k] = tableModel
	}
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		CreateInBatches(tableModels, len(tableModels)).Error
	if err != nil {
		return err
	}
	return nil
}

//CreateByMap 新建
func (m *CommonQuery) CreateByMap(ctx context.Context, tableName string, data map[string]interface{}) error {
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Table(tableName).
		Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *CommonQuery) DeleteByWhere(ctx context.Context, tableModel CommonTabler, conditions Conditions) error {
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Where(conditions.Where, conditions.Values...).
		Delete(tableModel).
		Error
	return err
}

//FindOne 根据主键单个查询
func (m *CommonQuery) FindOne(ctx context.Context, tableModel CommonTabler, value interface{}, hasCache ...bool) (err error) {
	//做异常检测
	defer m.checkAble()
	rawCache := false
	if len(hasCache) > 0 {
		rawCache = hasCache[0] && !m.CacheSignal.CacheOff
	}
	m.cacheQuery = rawCache
	pkName := tableModel.GetPKColumnName()
	m.FindOneCache(fmt.Sprintf(CacheTableRowKey, tableModel.TableName(), pkName, value), tableModel, func() {
		err = m.GormDb.Session(&gorm.Session{
			NewDB: true,
		}).WithContext(ctx).
			Where(fmt.Sprintf("%s=?", pkName), value).
			First(tableModel).Error

		return
	})

	return
}

//FindOneBy 单个查询
func (m *CommonQuery) FindOneBy(ctx context.Context, tableModel CommonTabler, column string, value interface{}) error {
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Where(fmt.Sprintf("`%s`=?", column), value).
		First(tableModel).Error
	if err != nil {
		return err
	}
	return nil
}

//FindOneBy 单个查询
func (m *CommonQuery) FindOneWithCondition(ctx context.Context, tableModel CommonTabler, conditions Conditions) error {
	err := m.GormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Select(conditions.Columns).
		Where(conditions.Where, conditions.Values...).
		Order(conditions.OrderBy).
		First(tableModel).Error
	if err != nil {
		return err
	}
	return nil
}

//Save 更新或新建
func (m *CommonQuery) Save(ctx context.Context, tableModel CommonTabler, hasCache ...bool) error {
	defer m.checkAble()
	tableModel.SetTime()
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).Save(tableModel).Error
	if err != nil {
		return err
	}

	rawCache := false
	if len(hasCache) > 0 {
		rawCache = hasCache[0]
	}
	if rawCache {
		pkName := tableModel.GetPKColumnName()
		m.UpdateCache(fmt.Sprintf(CacheTableRowKey, tableModel.TableName(), pkName, tableModel.GetPkValue()), tableModel)
	}
	return nil
}

//Updates 根据model多字段更新
func (m *CommonQuery) Updates(ctx context.Context, tableModel CommonTabler, updateMap map[string]interface{}) error {
	defer m.checkAble()
	selectSlice := make([]string, 0)
	pkColumn := strings.Trim(tableModel.GetPKColumnName(), "`")
	for k, _ := range updateMap {
		if k != pkColumn {
			selectSlice = append(selectSlice, k)
		}
	}
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Table(tableModel.TableName()).
		Where(fmt.Sprintf("%s=?", tableModel.GetPKColumnName()), tableModel.GetPkValue()).
		Select(selectSlice).
		Updates(updateMap).
		Error
	if err != nil {
		return err
	}
	return nil
}

//UpdatesByConditions 根据条件修改
func (m *CommonQuery) UpdatesByConditions(ctx context.Context, tableName string, conditions Conditions, updateMap map[string]interface{}) error {
	defer m.checkAble()

	selectSlice := make([]string, 0)
	for k, _ := range updateMap {
		selectSlice = append(selectSlice, k)
	}
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Table(tableName).
		Where(conditions.Where, conditions.Values...).
		Select(selectSlice).
		Updates(updateMap).
		Error
	if err != nil {
		return err
	}
	return nil
}

//Update 单字段更新
func (m *CommonQuery) Update(ctx context.Context, tableModel CommonTabler, column string, value interface{}) error {
	defer m.checkAble()
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Table(tableModel.TableName()).
		Where(fmt.Sprintf("%s=?", tableModel.GetPKColumnName()), tableModel.GetPkValue()).
		Update(column, value).
		Error
	if err != nil {
		return err
	}
	return nil
}

//Delete 删除
func (m *CommonQuery) Delete(ctx context.Context, tableModel CommonTabler, hasCache ...bool) error {
	defer m.checkAble()
	if strings.Trim(tableModel.GetPKColumnName(), "`") == PK {
		err := m.WriteGormDb.Session(&gorm.Session{
			NewDB: true,
		}).WithContext(ctx).Delete(tableModel).
			Error
		if err != nil {
			return err
		}
	} else {
		err := m.WriteGormDb.Session(&gorm.Session{
			NewDB: true,
		}).WithContext(ctx).
			Where(fmt.Sprintf("%s=?", tableModel.GetPKColumnName()), tableModel.GetPkValue()).
			Delete(tableModel).Error
		if err != nil {
			return err
		}
	}
	rawCache := false
	if len(hasCache) > 0 {
		rawCache = hasCache[0]
	}
	if rawCache {
		pkName := tableModel.GetPKColumnName()
		m.DeleteCache(fmt.Sprintf(CacheTableRowKey, tableModel.TableName(), pkName, tableModel.GetPkValue()))
	}
	return nil
}

//DeleteByPk 根据主键删除
func (m *CommonQuery) DeleteByPk(ctx context.Context, tableModel CommonTabler, id interface{}) error {
	defer m.checkAble()
	err := m.WriteGormDb.Session(&gorm.Session{
		NewDB: true,
	}).WithContext(ctx).
		Where(fmt.Sprintf("%s=?", tableModel.GetPKColumnName()), id).
		Delete(tableModel).
		Error
	if err != nil {
		return err
	}
	return nil
}

//FindOneCache 查询缓存
func (m *CommonQuery) FindOneCache(cacheKey string, res interface{}, handler func()) {
	if m.cacheQuery {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.CacheRedis.GetCtx(ctx, cacheKey, res)
		if err != nil {
			handler()
			if m.CacheRedis.IsNotFound(err) {
				_ = m.CacheRedis.SetCtx(ctx, cacheKey, res)
			} else {
				m.CacheSignal.CacheErr = err
			}
		}

		return
	}
	handler()

	return
}

// QueryCache 查询有限缓存
func (m *CommonQuery) QueryCache(ctx context.Context, querySql string, values []interface{}, res interface{}, cacheTime time.Duration) error {
	if !m.CacheSignal.CacheOff {
		cacheKey := querySql
		for _, v := range values {
			cacheKey = strings.Replace(cacheKey, "?", fmt.Sprintf("'%v'", v), 1)
		}
		cacheKey = CachePrefix + utils.Md5Encode(cacheKey)
		cacheCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.CacheRedis.GetCtx(cacheCtx, cacheKey, res)
		if err != nil {
			err2 := m.Query(ctx, querySql, values, res)
			if err2 != nil {
				return err2
			}
			if m.CacheRedis.IsNotFound(err) {
				_ = m.CacheRedis.SetWithExpireCtx(cacheCtx, cacheKey, res, cacheTime)
			} else {
				m.CacheSignal.CacheErr = err
			}
		}

		return nil
	}

	err := m.Query(ctx, querySql, values, res)
	if err != nil {
		return err
	}

	return nil
}

//Query 查询
func (m *CommonQuery) Query(ctx context.Context, querySql string, values []interface{}, res interface{}) error {
	err := m.GormDb.Session(&gorm.Session{NewDB: true}).WithContext(ctx).Raw(querySql, values...).Find(res).Error
	return err
}

func (m *CommonQuery) Exec(ctx context.Context, querySql string, values ...interface{}) error {
	err := m.WriteGormDb.Session(&gorm.Session{NewDB: true}).WithContext(ctx).Exec(querySql, values...).Error
	return err
}

//UpdateCache 更新缓存
func (m *CommonQuery) UpdateCache(cacheKey string, res interface{}) {
	if m.CacheSignal.CacheOff {
		m.WriteErr(cacheKey)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.CacheRedis.SetCtx(ctx, cacheKey, res)
	if err != nil && !m.CacheRedis.IsNotFound(err) {
		m.CacheSignal.CacheErr = err
		m.WriteErr(cacheKey)
	}

	return
}

//DeleteCache 删除缓存
func (m *CommonQuery) DeleteCache(cacheKey string) {
	if m.CacheSignal.CacheOff {
		m.WriteErr(cacheKey)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.CacheRedis.DelCtx(ctx, cacheKey)
	if err != nil && !m.CacheRedis.IsNotFound(err) {
		m.CacheSignal.CacheErr = err
		m.WriteErr(cacheKey)
	}
	return
}

//redis 异常检测
func (m *CommonQuery) checkAble() {
	if !m.cacheQuery {
		return
	}
	if m.CacheSignal.CacheErr == nil {
		return
	}
	m.CacheSignal.CacheChan <- commonconst.CacheOffSignal
	//置为断联，并重置错误
	m.CacheSignal.CacheOff = true
	m.CacheSignal.CacheErr = nil
	return
}

func (m *CommonQuery) WriteErr(cacheKey string) {
	if len(m.CacheSignal.CacheErrFile) == 0 {
		return
	}
	dirPath, _ := os.Getwd()
	dirPath += commonconst.CacheErrorLogName
	_, e := os.Stat(dirPath)
	if e != nil {
		if os.IsNotExist(e) {
			if e := os.MkdirAll(dirPath, os.ModePerm); e != nil {
				fmt.Println(fmt.Sprintf("%v\n%s", e, debug.Stack()))
				return
			}
		} else {
			return
		}
	}
	fileName := dirPath + m.CacheSignal.CacheErrFile
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logx.Errorf("cache write error's key open err:%v", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(cacheKey + "\n")
	if err != nil {
		logx.Errorf("cache write error's key encode err:%v", err)
		return
	}
}
