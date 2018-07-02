package model

/**
该文件用来放置模型的定义
*/
import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"media_framwork/tool"
	"strconv"
	"time"
)

const (
	COMMONUSER = 1
	ADMIN      = 2
)

type Model struct {
	ID        int `gorm:"primary_key" `
	CreatedAt time.Time
	UpdatedAt time.Time
	//1 undelete 2 delete
	SoftDelete int `gorm:"default:1"`
}

type User struct {
	Model
	Username string `gorm:"type:char(64);not null" `
	Password string `gorm:"type:char(64);not null" `
	Nickname string `gorm:"type:char(64);not null" `
	Avatar   string
	//1 COMMON 2 ADMIN
	Authority   int `gorm:"default:1"`
	UserRecords []*UserRecord
}

type Media struct {
	Model
	Title           string `gorm:"type:varchar(255);not null;index"`
	Introduction    string `gorm:"type:varchar(1000);not null"`
	Cover           string
	Categories      []*Category `gorm:"many2many:media_categories"`
	MediaAttributes []*MediaAttribute
	StarCount       int `gorm:"-"`
}

type MediaAttribute struct {
	ID             int `gorm:"primary_key"`
	Media          *Media
	MediaID        int `gorm:"index;not null"`
	Uri            string
	Filename       string
	Description    string
	DownloadStatus int `gorm:"type:tinyint;default 0"`
}

type Category struct {
	Model
	Name   string   `gorm:"type:varchar(100);not null;index"`
	Medias []*Media `gorm:"many2many:media_categories"`
}

type UserRecord struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
	UserID    int `gorm:"index;not null"`
	Media     *Media
	MediaID   int `gorm:"index;not null"`
}

type Star struct {
	Model
	User    *User
	UserID  int `gorm:"index;not null"`
	Media   *Media
	MediaID int `gorm:"index;not null"`
}

type Comment struct {
	Model
	User            *User
	UserID          int      `gorm:"index;not null"`
	ParentComment   *Comment
	ParentCommentID int      `gorm:"index" json:"Parent_ID,string"`
	Media           *Media
	MediaID         int    `gorm:"index;not null" json:"Media_ID,string"`
	Content         string `gorm:"varchar(1000);not null" json:"Content"`
}

type Page struct {
	Limit    int
	CurPage  int
	MaxPage  int
	PrevPage int
	NextPage int
	Count    int
	CurCount int
	PrevLink string
	NextLink string
	c        *gin.Context
}

/**
创建指定限制和页数的page
*/
func NewPage(c *gin.Context, limit int) *Page {
	v := tool.GetInt(c.DefaultQuery("page", "1"))
	if v < 1 {
		v = 1
	}
	return &Page{
		Limit:   limit,
		CurPage: v,
		c:       c,
	}
}

/**
生成默认限制为15的page
*/
func DefaultPage(c *gin.Context) *Page {
	return NewPage(c, 10)
}

/**
查找数据并且生成分页信息
*/
func (p *Page) Find(i interface{}, db *gorm.DB, delete ...bool) {
	curDB := db.Model(i)
	if len(delete) == 1 {
		if delete[0] {
			curDB = curDB.Where("soft_delete=2")
		} else {
			curDB = curDB.Where("soft_delete=1")
		}
	}
	curDB.
		Count(&p.Count).
		Offset((p.CurPage - 1) * p.Limit).
		Limit(p.Limit).
		Find(i).
		Count(&p.CurCount)
	mod := p.Count % p.Limit
	if mod == 0 {
		p.MaxPage = p.Count / p.Limit
	} else {
		p.MaxPage = p.Count/p.Limit + 1
	}

	if p.CurPage+1 <= p.MaxPage {
		p.NextPage = p.CurPage + 1
	} else {
		p.NextPage = p.CurPage
	}
	if p.CurPage-1 >= 1 {
		p.PrevPage = p.CurPage - 1
	} else {
		p.PrevPage = p.CurPage
	}
	p.generateLink()
}

/**
生成分页链接
*/
func (p *Page) generateLink() {
	uri := p.c.Request.URL
	path := uri.Path

	raw := uri.Query()
	raw.Set("page", strconv.Itoa(p.NextPage))
	p.NextLink = path + "?" + raw.Encode()

	raw = uri.Query()
	raw.Set("page", strconv.Itoa(p.PrevPage))
	p.PrevLink = path + "?" + raw.Encode()

}

/**
批量删除
*/
func Delete(i interface{}, ids ...int) error {
	db.Model(i).Where("id in (?)", ids).Update("soft_delete", 2)
	return db.Error
}

/**
批量恢复
*/
func Recover(i interface{}, ids ...int) error {
	db.Model(i).Where("id in (?)", ids).Update("soft_delete", 1)
	return db.Error
}

/**
获取对象数量
*/
func Count(i interface{}) int {
	count := 0
	db.Model(i).Count(&count)
	return count
}

/**
通过ID找对象
*/
func FindById(i interface{}, id int) (interface{}, error) {
	count := 0
	db.Model(i).Where("id=?", id).Find(i).Count(&count)
	if count == 0 {
		return nil, errors.New("数据不存在")
	}
	return i, nil
}

func Exist(i interface{}, id int) bool {
	c := 0
	db.Model(i).Where("id=?", id).Count(&c)
	return c != 0
}
