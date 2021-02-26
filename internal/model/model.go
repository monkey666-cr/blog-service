package model

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"blog-service/global"
	"blog-service/pkg/setting"
)

type Model struct {
	ID         uint32         `gorm:"primary_key" json:"id"`
	CreateBy   string         `json:"created_by"`
	ModifiedBy string         `json:"modified_by"`
	CreatedAt  time.Time      `json:"created_on"`
	ModifiedAt time.Time      `json:"modified_on"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
			databaseSetting.UserName,
			databaseSetting.Password,
			databaseSetting.Host,
			databaseSetting.DBName,
			databaseSetting.Charset,
			databaseSetting.ParseTime,
		),
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode != "debug" {
		db.Logger.LogMode(logger.Info)
	}

	idb, _ := db.DB()
	idb.SetMaxIdleConns(databaseSetting.MaxIdleConnes)
	idb.SetMaxOpenConns(databaseSetting.MaxOpenConnes)

	return db, nil
}

func MysqlTables(db *gorm.DB) error {
	if err := db.AutoMigrate(
		Article{},
	); err != nil {
		return err
	}
	global.Logger.Info(nil, "register table success")
	return nil
}
