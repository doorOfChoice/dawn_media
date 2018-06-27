package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"media_framwork/model"
)

func MediaCreate(media *model.Media, ms []*model.MediaSharpness, cs []*model.Category) (*model.Media, error) {
	db := model.DB()
	tx := db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()

		}
	}()
	tx.Create(media)
	if err := tx.Error; err != nil {
		return nil, err
	}
	if err := mediaAddRef(media, ms, cs, tx); err != nil {
		return nil, err
	}
	return media, tx.Commit().Error
}

func mediaAddRef(media *model.Media, ms []*model.MediaSharpness, cs []*model.Category, tx *gorm.DB) error {
	for _, m := range ms {
		m.Media = media
		if err := tx.Save(m).Error; err != nil {
			return err
		}
	}
	tx.
		Model(media).
		Association("Categories").
		Append(cs)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func MediaUpdate(
	id int,
	ms []*model.MediaSharpness,
	cs []*model.Category,
) (*model.Media, error) {
	db := model.DB()
	tx := db.Begin()
	var media model.Media
	var count int
	media.ID = id
	db.Find(&media).Count(&count)
	if count == 0 {
		tx.Rollback()
		return nil, errors.New("媒体文件不存在")
	}
	tx.Model(&media).Association("Categories").Clear()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := mediaAddRef(&media, ms, cs, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &media, tx.Commit().Error
}
