package model

import (
	"errors"
	"media_framwork/conf"
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
	if err := tx.Save(media).Error; err != nil {
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
		Joins(`JOIN media_categories as mc ON mc.category_id=id JOIN media ON mc.media_id=media.id and media.id=?`, m.ID).
		Find(&m.Categories)
	db.Model(m).
		Association("MediaAttributes").
		Find(&m.MediaAttributes)
	m.Cover = conf.C().CoverMap + m.Cover
	for _, ma := range m.MediaAttributes {
		ma.Uri = conf.C().MediaMap + ma.Uri
	}
	return db.Error
}

