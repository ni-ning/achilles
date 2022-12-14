package model

import (
	"achilles/pkg/setting"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 表名小写复数 users，字段单词小写下划线，即蛇形Python风格
type User struct {
	ID        uint64    `gorm:"primarykey;" json:"id"`
	Username  string    `gorm:"index" json:"username"`
	Age       uint8     `gorm:"default:18" json:"age"`  // 创建时传入0时，保存会为18，有点奇怪
	Age2      *uint8    `gorm:"default:18" json:"age2"` // 保证0、''、false 等零值会写入到数据库
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsDeleted bool      `json:"-"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	dsn := fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// TODO 再研究一下
	// if global.ServerSetting.RunMode == "debug" {
	// 	db.LogMode(true)
	// }
	// db.SingularTable(true)
	// db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	// db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}
