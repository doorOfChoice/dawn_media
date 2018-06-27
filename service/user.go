package service

import (
	"errors"
	"media_framwork/model"
)
/**
 * 在数据库中创建用户
 */
func UserCreate(user *model.User) (*model.User, error) {
	var count int
	model.DB().Model(&model.User{}).Where("username=?", user.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户已被注册")
	}
	model.DB().Save(&user)
	if err := model.DB().Error; err != nil {
		return nil, err
	}
	return user, nil
}
