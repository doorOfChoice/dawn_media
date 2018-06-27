package model

/**
该文件用来放置模型的定义
*/
import (
	"time"
)

type Model struct {
	ID         int       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"createAt"`
	UpdatedAt  time.Time `json:"updateAt"`
	SoftDelete int       `gorm:"default:0"`
}

type User struct {
	Model
	Username    string        `gorm:"type:char(64);not null" json:"username"`
	Password    string        `gorm:"type:char(64);not null" json:"password"`
	Sex         int           `gorm:"type:tinyint;default 0" json:"sex"`
	UserRecords []*UserRecord `json:"userRecords"`
}

type Media struct {
	Model
	Title            string            `gorm:"type:varchar(255);not null;index" json:"title"`
	Introduction     string            `gorm:"type:varchar(1000);not null" json:"introduction"`
	DownloadState    int8              `gorm:"type:tinyint;default 0" json:"downloadState"`
	Categories       []*Category       `gorm:"many2many:media_categories" json:"categories"`
	MediaSharpnesses []*MediaSharpness `json:"mediaSharpnesses"`
}

type Category struct {
	Model
	Name   string   `gorm:"type:varchar(100);not null;index" json:"name"`
	Medias []*Media `gorm:"many2many:media_categories" json:"medias"`
}

type Sharpness struct {
	Model
	Name string `gorm:"type:char(30);not null;index" json:"name"`
}

type UserRecord struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      *User     `json:"user"`
	UserID    int       `gorm:"index" json:"-"`
	Media     *Media    `json:"media"`
	MediaID   int       `gorm:"index" json:"-"`
}

type Star struct {
	Model
	User    *User  `json:"user"`
	UserID  int    `gorm:"index" json:"-"`
	Media   *Media `json:"media"`
	MediaID int    `gorm:"index"json:"-"`
}

type MediaSharpness struct {
	ID          int        `gorm:"primary_key;auto_increment" json:"id"`
	Media       *Media     `json:"media"`
	MediaID     int        `gorm:"index" json:"-"`
	Sharpness   *Sharpness `json:"sharpness"`
	SharpnessId int        `gorm:"index" json:"-"`
	Uri         string     `json:"uri"`
}
