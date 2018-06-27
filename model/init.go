package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Init() {
	if db != nil {
		return
	}
	t, err := gorm.Open("mysql", "root:1997@/media_web?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	t.AutoMigrate(
		&User{},
		&UserRecord{},
		&Media{},
		&Category{},
		&Sharpness{},
		&Star{},
		&MediaSharpness{},
	)

	db = t
}

func DB() *gorm.DB {
	if db == nil {
		panic("gorm db is not exist")
	}
	return db
}
