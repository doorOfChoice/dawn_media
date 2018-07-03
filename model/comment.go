package model

import (
	"errors"
	"github.com/jinzhu/gorm"
)

/**
创建评论
 */
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

/**
获取评论
 */
func (c *Comment) Get() error {
	return db.
		Preload("User").
		Preload("Media").
		Preload("ParentComment.User").
		Find(c).Error
}
/**
获取具有条件筛选的comment集合
mid 为媒体id，0为不查找
uid 为用户id，0为不查找
curTime 为某一时间段钱
 */
func CommentConditionDB(mid, uid int, curTime string) *gorm.DB {
	//后面Media验证一下是否存在
	curDB := db.
		Preload("User").
		Preload("ParentComment.User").
		Preload("Media")
	if mid != 0 {
		curDB = curDB.Where("media_id=?", mid)
	}
	if uid != 0 {
		curDB = curDB.Where("user_id=?", uid)
	}
	return curDB.
		Where(" unix_timestamp(created_at)<?", curTime).
		Order("created_at desc")
}


/**
获取主页展示的评论
默认限制4条
 */
func GetIndexNewComments() []*Comment {
	comments := make([]*Comment, 0)
	db.
		Preload("User").
		Preload("Media").
		Order("created_at desc").
		Limit(4).
		Find(&comments)
	return comments
}
