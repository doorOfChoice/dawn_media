package model

import "github.com/jinzhu/gorm"

/**
获取用户记录的DB
 */
func UserRecordDB(id int) *gorm.DB {
	return db.Preload("Media").Order("updated_at desc").Where("user_id=?", id)
}

/**
更新浏览记录
如果有某一个浏览记录，只是更新时间
 */
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
	return db.Model(record).Update(record).Error
}
