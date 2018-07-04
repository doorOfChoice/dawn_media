package model

import (
	"errors"
	"dawn_media/conf"
	"dawn_media/tool"
)

/**
创建用户
第一个用户为管理员
*/
func (u *User) Create() error {
	count := 0
	db.
		Model(u).
		Where("username=?", u.Username).
		Count(&count)
	if count > 0 {
		return errors.New("用户已经被注册")
	}

	db.Model(u).Count(&count)
	if count == 0 {
		u.Authority = ADMIN
	}
	u.Password = tool.Md5EncodeWithSalt(u.Password, conf.C().PassSalt)
	return db.Save(u).Error
}

/**
更新用户
*/
func (u *User) Update() error {
	user := &User{}
	user.ID = u.ID

	if err := user.Get(); err != nil {
		return errors.New("用户不存在")
	}
	if u.Password != "" {
		user.Password = tool.Md5EncodeWithSalt(u.Password, conf.C().PassSalt)
	}
	user.Nickname = u.Nickname
	user.Authority = u.Authority
	user.Avatar = u.Avatar
	return db.Model(user).Update(user).Error
}

/**
获取用户详细信息
*/
func (u *User) Get() error {
	count := 0
	db.Find(u).Count(&count)
	if count == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

/**
获取管理员主页上最新注册的用户
*/
func GetAdminIndexNewUser() []*User {
	users := make([]*User, 0)
	db.
		Where("soft_delete=1").
		Order("created_at desc").
		Limit(5).
		Find(&users)
	return users
}

/**
通过用户名和密码获取用户
*/
func FindUserByLogin(username, password string) (*User, error) {
	count := 0
	user := &User{}
	password = tool.Md5EncodeWithSalt(password, conf.C().PassSalt)
	db.
		Where("username=? and password=?", username, password).
		Find(user).
		Count(&count)
	if count == 0 {
		return nil, errors.New("用户名或者密码错误")
	}
	return user, nil
}
