package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Init(params ...interface{}) {
	if db != nil {
		return
	}
	t, err := gorm.Open("mysql", "root:1997@/media_web?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	t.LogMode(true)
	t.AutoMigrate(
		&User{},
		&UserRecord{},
		&Media{},
		&Category{},
		&Star{},
		&Comment{},
		&MediaAttribute{},
	)
	t.Table("media_categories").AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")
	t.Table("media_categories").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("user_records").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("user_records").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	t.Table("stars").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("stars").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	t.Table("comments").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("comments").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	t.Table("media_attributes").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	db = t
}

func DB() *gorm.DB {
	if db == nil {
		panic("gorm db is not exist")
	}
	return db
}
