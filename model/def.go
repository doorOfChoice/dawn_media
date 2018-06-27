package model
/**
	该文件用来放置模型的定义
 */
import (
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createAt"`
	UpdatedAt time.Time  `json:"updateAt"`
	DeletedAt *time.Time `gorm:"deleteAt"`
}

type User struct {
	*Model
	Username    string        `gorm:"type:char(64);not null" json:"username"`
	Password    string        `gorm:"type:char(64);not null" json:"password"`
	Sex         int           `gorm:"type:tinyint;default 0" json:"sex"`
	UserRecords []*UserRecord `json:"userRecords"`
}

type Media struct {
	*Model
	Title            string      `gorm:"type:varchar(255);not null;index"`
	Introduction     string      `gorm:"type:varchar(1000);not null"`
	DownloadState    int8        `gorm:"type:tinyint;default 0"`
	Categories       []*Category `gorm:"many2many:media_categories"`
	MediaSharpnesses []*MediaSharpness
}

type Category struct {
	*Model
	Name   string   `gorm:"type:varchar(100);not null;index"`
	Medias []*Media `gorm:"many2many:media_categories"`
}

type Sharpness struct {
	*Model
	Name string `gorm:"type:char(30);not null;index"`
}

type UserRecord struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
	UserID    uint
	Media     *Media
	MediaID   uint
}

type Star struct {
	*Model
	User    *User
	UserID  uint
	Media   *Media
	MediaID uint
}

type MediaSharpness struct {
	ID          uint `gorm:"primary_key;auto_increment"`
	Media       *Media
	MediaID     uint
	Sharpness   *Sharpness
	SharpnessId uint
	Uri         string
}
