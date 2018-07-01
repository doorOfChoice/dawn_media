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
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updateAt"`
	//1 undelete 2 delete
	SoftDelete int `gorm:"default:1"`
}

type User struct {
	Model
	Username string `gorm:"type:char(64);not null" json:"username"`
	Password string `gorm:"type:char(64);not null" json:"password"`
	Nickname string `gorm:"type:char(64);not null" json:"nickname"`
	Avatar   string
	//1 COMMON 2 ADMIN
	Authority   int           `gorm:"default:1"`
	UserRecords []*UserRecord `json:"userRecords"`
}

type Media struct {
	Model
	Title           string `gorm:"type:varchar(255);not null;index" json:"title"`
	Introduction    string `gorm:"type:varchar(1000);not null" json:"introduction"`
	Cover           string
	Categories      []*Category `gorm:"many2many:media_categories" json:"categories"`
	MediaAttributes []*MediaAttribute
}

type MediaAttribute struct {
	ID             int `gorm:"primary_key" json:"id"`
	Media          *Media
	MediaID        int
	Uri            string
	Filename       string
	Description    string
	DownloadStatus int `gorm:"type:tinyint;default 0" json:"downloadState"`
}

type Category struct {
	Model
	Name   string   `gorm:"type:varchar(100);not null;index" json:"name"`
	Medias []*Media `gorm:"many2many:media_categories" json:"medias"`
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

type Comment struct {
	Model
	User    *User
	UserID  int `gorm:"index"`
	Media   *Media
	MediaID int    `gorm:"index"`
	Content string `gorm:"varchar(1000);not null"`
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
	count := 0
	curDB := db.Model(i)
	if len(delete) == 1 {
		if delete[0] {
			curDB = curDB.Where("soft_delete=2")
		} else {
			curDB = curDB.Where("soft_delete=1")
		}
	}
	curDB.
		Count(&count).
		Offset((p.CurPage - 1) * p.Limit).
		Limit(p.Limit).
		Find(i).
		Count(&p.CurCount)
	p.Count = count
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



