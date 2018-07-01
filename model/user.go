package model

import (
	"errors"
	"media_framwork/tool"
	"media_framwork/conf"
)

func (u *User) Create() error {
	count := 0
	db.
		Model(u).
		Where("username=?", u.Username).
		Count(&count)
	if count > 0 {
		return errors.New("用户已经被注册")
	}
	u.Password = tool.Md5EncodeWithSalt(u.Password, conf.C().PassSalt)
	return db.Save(u).Error
}

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

func (u *User) Get() error {
	count := 0
	db.Find(u).Count(&count)
	if count == 0 {
		return errors.New("用户不存在")
	}
	u.Avatar = conf.C().AvatarMap + u.Avatar
	return nil
}

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
