package model

import "errors"

/**
点赞
没有点赞则点赞
点过赞则取消点赞

返回值 ： 是否点赞，当前点赞总数，错误
 */
func ToggleStar(uid, mid int) (bool, int, error) {
	var err error
	create := false
	count := 0
	star := &Star{}
	//查找是否点赞
	err = db.
		Model(&Star{}).Where("user_id=? and media_id=?", uid, mid).
		First(star).
		Count(&count).Error
	if !Exist(&Media{}, mid) || !Exist(&User{}, uid) {
		return false, 0, errors.New("媒体或者用户不存在")
	}
	//点赞
	if count == 0 {
		create = true
		star.UserID = uid
		star.MediaID = mid
		err = db.Save(star).Error
	//取消赞
	} else {
		create = false
		err = db.Delete(star).Error
	}
	//db.Model查询全文的对象时候，尽量不要重复用之前的对象作为参数
	db.Model(&Star{}).Where("media_id=?", mid).Count(&count)
	return create, count, err
}
/**
判断某用户是否赞过某个影像
mid 媒体id
uid 用户id
 */
func HasStared(mid, uid int) bool {
	count := 0
	db.
		Model(&Star{}).Where("user_id=? and media_id=?", uid, mid).
		Count(&count)
	return count != 0
}
