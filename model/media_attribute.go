package model

import (
	"github.com/jinzhu/gorm"
	"errors"
	"strconv"
)

func (ma *MediaAttribute) Update(i ...interface{}) error {
	var tx *gorm.DB
	if len(i) != 0 {
		tx = i[0].(*gorm.DB)
	}
	count := 0
	tx.Model(ma).Count(&count)
	if count == 0 {
		return errors.New("属性" + strconv.Itoa(ma.ID) + "不存在")
	}
	return tx.Model(ma).Update(ma).Error
}

func (ma *MediaAttribute) Create(i ...interface{}) error {
	var tx *gorm.DB
	if len(i) != 0 {
		tx = i[0].(*gorm.DB)
	}
	return tx.Save(ma).Error
}

