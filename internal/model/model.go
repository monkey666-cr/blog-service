package model

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"blog-service/global"
	"blog-service/pkg/setting"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreateBy   string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	fileName := global.AppSetting.LogSavePath + "/" + "sql" + global.AppSetting.LogFileExt

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
			databaseSetting.UserName,
			databaseSetting.Password,
			databaseSetting.Host,
			databaseSetting.DBName,
			databaseSetting.Charset,
			databaseSetting.ParseTime,
		),
	}), &gorm.Config{
		Logger: logger.New(log.New(&lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    500,
			MaxAge:     10,
			MaxBackups: 1024,
			LocalTime:  true,
		}, "", log.LstdFlags), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel: func() logger.LogLevel {
				if global.ServerSetting.RunMode == "debug" {
					return logger.Silent
				}
				return logger.Info
			}(),
			Colorful: false,
		}),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
