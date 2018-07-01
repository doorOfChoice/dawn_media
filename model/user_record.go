package model

import (
	"github.com/jinzhu/gorm"
)

func UserRecordDB(id int) *gorm.DB {
	return db.Preload("Media").Order("updated_at desc").Where("user_id=?", id)
}

func UserRecordUpdate(uid int, mid int) error {
	record := &UserRecord{MediaID: mid, UserID: uid}
	count := 0
	db.
		Where("user_id=? and media_id=?", uid, mid).
		Find(record).
		Count(&count)
	if count == 0 {
		return db.Save(record).Error
	}
	return db.Update(record).Error
}
