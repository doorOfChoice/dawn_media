package model

import (
	"errors"
	"github.com/jinzhu/gorm"
)

func (c *Comment) Create() error {
	if c.ParentCommentID != 0 && !Exist(c, c.ParentCommentID) {
		return errors.New("目标评论不存在")
	}

	media := &Media{}
	media.ID = c.MediaID
	if media.ID == 0 || !Exist(media, media.ID) {
		return errors.New("目标媒体不存在")
	}

	return db.Save(c).Error
}

func (c *Comment) Get() error {
	return db.
		Preload("User").
		Preload("Media").
		Preload("ParentComment.User").
		Find(c).Error
}

func CommentConditionDB(mid int, curTime string) *gorm.DB {
	//后面Media验证一下是否存在
	return db.
		Preload("User").
		Preload("ParentComment.User").
		Where("media_id=?", mid).
		Where(" unix_timestamp(created_at)<?", curTime).
		Order("created_at desc")
}

func GetIndexNewComments()[]*Comment {
	comments := make([]*Comment, 0)
	db.
		Preload("User").
		Preload("Media").
		Order("created_at desc").
		Limit(4).
		Find(&comments)
	return comments
}