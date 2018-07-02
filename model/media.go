package model

import (
	"errors"
	"github.com/jinzhu/gorm"
)

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
	for _, ma := range m.MediaAttributes {
		if err := ma.Create(tx); err != nil {
			tx.Rollback()
			return err
		}
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
	media.Cover = m.Cover
	tx := db.Begin()
	if err := tx.Model(media).Update(media).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(media).Association("Categories").Clear().Append(m.Categories)
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, ma := range m.MediaAttributes {
		if err := ma.Update(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (m *Media) Get() error {
	m.Categories = make([]*Category, 0)
	m.MediaAttributes = make([]*MediaAttribute, 0)
	count := 0
	db.Find(m).Count(&count)
	if count == 0 {
		return errors.New("没有找到媒体文件")
	}
	db.
		Joins(`JOIN media_categories as mc ON mc.category_id=id`).
		Joins(`JOIN media ON mc.media_id=media.id and media.id=?`, m.ID).
		Where("categories.soft_delete=1").
		Find(&m.Categories)
	db.Model(m).
		Association("MediaAttributes").
		Find(&m.MediaAttributes)
	return db.Error
}

func GetIndexRandomMedia() []*Media {
	medias := make([]*Media, 0)
	db.Where("soft_delete=1").Order("RAND()").Limit(6).Find(&medias)
	return medias
}

func GetIndexHotMedia() []*Media {
	medias := make([]*Media, 0)
	db.
		Model(&Media{}).
		Select("media.*, count(stars.media_id) as star_count").
		Joins("left join stars on stars.media_id=media.id").
		Where("media.soft_delete=1").
		Group("media.id").
		Order("star_count desc").
		Find(&medias)
	return medias
}

func GetIndexNewMedia() []*Media {
	medias := make([]*Media, 0)
	db.Where("soft_delete=1").Order("created_at desc").Limit(3).Find(&medias)
	return medias
}

// byTime, byHot, 0 无 1 高到低 2 低到高
func ConditionSelectMediaDB(byTime, byHot, category int, name string) *gorm.DB {
	curDB := db
	curDB = curDB.
		Select("media.*, count(stars.media_id) as star_count").
		Joins("left join stars on stars.media_id=media.id")
	if category != 0 {
		curDB = curDB.
			Joins("left join media_categories as mc on mc.media_id=media.id").
			Joins("left join categories on mc.category_id=categories.id").
			Where("categories.id=?", category)

	}
	curDB = curDB.Group("media.id")
	if name != "" {
		curDB = curDB.Where("title like ?", "%"+name+"%")
	}
	if byTime == 1 {
		curDB = curDB.Order("created_at desc")
	} else if byTime == 2 {
		curDB = curDB.Order("created_at asc")
	}
	if byHot == 1 {
		curDB = curDB.Order("star_count desc")
	} else if byHot == 2 {
		curDB = curDB.Order("star_count asc")
	}
	return curDB
}
