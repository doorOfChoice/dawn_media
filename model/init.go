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
		&Sharpness{},
		&Star{},
		&MediaSharpness{},
	)
	t.Table("media_categories").AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")
	t.Table("media_categories").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("media_sharpnesses").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("media_sharpnesses").AddForeignKey("sharpness_id", "sharpnesses(id)", "CASCADE", "CASCADE")
	t.Table("user_records").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("user_records").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	t.Table("stars").AddForeignKey("media_id", "media(id)", "CASCADE", "CASCADE")
	t.Table("stars").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	if len(params) == 1 && params[0].(bool) {
		t.Create(&Sharpness{Model:Model{ID:1}, Name:"标清"})
		t.Create(&Sharpness{Model:Model{ID:2}, Name:"高清"})
		t.Create(&Sharpness{Model:Model{ID:3}, Name:"标清"})
		t.Create(&Category{Model:Model{ID:3}, Name:"FK"})
	}
	db = t
}

func DB() *gorm.DB {
	if db == nil {
		panic("gorm db is not exist")
	}
	return db
}
