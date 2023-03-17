package gormpool

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type GormConfig struct {
	DataSource      string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	LogLevel        logger.LogLevel
}

func CoonMysql(gormConfig GormConfig) (*gorm.DB, error) {
	//启动Gorm支持
	db, err := gorm.Open(mysql.Open(gormConfig.DataSource), &gorm.Config{
		Logger: logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 1 * time.Second,
			LogLevel:      gormConfig.LogLevel,
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项
		},
	})
	//如果出错就GameOver了
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(gormConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(gormConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(gormConfig.ConnMaxLifetime) * time.Second)
	return db, nil
}
