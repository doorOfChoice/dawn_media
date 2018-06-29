package model

/**
该文件用来放置模型的定义
*/
import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"media_framwork/tool"
	"strconv"
	"time"
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
	Username    string        `gorm:"type:char(64);not null" json:"username"`
	Password    string        `gorm:"type:char(64);not null" json:"password"`
	Sex         int           `gorm:"type:tinyint;default 0" json:"sex"`
	UserRecords []*UserRecord `json:"userRecords"`
}

type Media struct {
	Model
	Title         string `gorm:"type:varchar(255);not null;index" json:"title"`
	Introduction  string `gorm:"type:varchar(1000);not null" json:"introduction"`
	S360          string
	S480          string
	S720          string
	S1080         string
	DownloadState int8        `gorm:"type:tinyint;default 0" json:"downloadState"`
	Categories    []*Category `gorm:"many2many:media_categories" json:"categories"`
}

type MediaAttribute struct {
	ID             int `gorm:"primary_key" json:"id"`
	Media          *Media
	MediaID        int
	Uri            string
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
func (p *Page) Find(i interface{}, db *gorm.DB) {
	count := 0
	db.Model(i).
		Count(&count).
		Offset((p.CurPage - 1) * p.Limit).
		Limit(p.Limit).
		Find(i)
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

	log.Println(p.PrevPage, p.NextPage, p.PrevLink, p.NextLink)
}

/**
已经被删除的DB
*/
func DeleteDB(i interface{}) *gorm.DB {
	return db.Model(i).Where("soft_delete=2")
}

/**
未被删除的DB
*/
func UnDeleteDB(i interface{}) *gorm.DB {
	return db.Model(i).Where("soft_delete=1")
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

/**
创建分类
*/
func (c *Category) Create() error {
	count := 0
	db.Model(Category{}).Where("name=?", c.Name).Count(&count)
	if count != 0 {
		return errors.New("分类已经存在")
	}
	db.Save(c)
	return db.Error
}

/**
更新分类
*/
func (c *Category) Update() error {
	if u, err := FindById(&Category{}, c.ID); err != nil {
		return err
	} else {
		ct := u.(*Category)
		count := 0
		db.Model(c).Where("name=?", c.Name).Count(&count)
		if count != 0 {
			return errors.New("标签名已被占用")
		}
		ct.Name = c.Name
		db.Save(ct)
		return db.Error
	}
}

func (m *Media) Create() error {
	tx := db.Begin()
	tx.Save(m)
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(m).Association("Categories").Append(m.Categories)
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (m *Media) Update() error {
	var media *Media
	if t, err := FindById(&Media{}, m.ID); err != nil {
		return err
	} else {
		media = t.(*Media)
	}
	media.Title = m.Title
	media.Introduction = m.Introduction
	media.S360 = m.S360
	media.S480 = m.S480
	media.S720 = m.S720
	media.S1080 = m.S1080
	tx := db.Begin()
	if err := tx.Save(media).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(media).Association("Categories").Clear().Append(m.Categories)
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (m *Media) Get() error {
	m.Categories = make([]*Category, 0)
	count := 0
	db.Find(m).Count(&count)
	if count == 0 {
		return errors.New("没有找到媒体文件")
	}
	db.
		Joins(`JOIN media_categories as mc ON mc.category_id=id JOIN media ON mc.media_id=media.id`).
		Find(&m.Categories)
	return db.Error
}
