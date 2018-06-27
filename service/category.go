package service

import (
	"errors"
	"media_framwork/model"
)

/**
标签创建
*/
func CategoryCreate(c *model.Category) (*model.Category, error) {
	db := model.DB()
	count := 0
	db.Model(&model.Category{}).Where("name=?", c.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("分类已经存在")
	}
	db.Save(c)
	if db.Error != nil {
		return nil, db.Error
	}
	return c, nil
}

/**
标签更新基本信息
*/
func CategoryBaseUpdate(id int, c *model.Category) (*model.Category, error) {
	db := model.DB()
	cg := model.Category{Model: model.Model{}}
	cg.ID = id
	count := 0
	//查看标签是否存在
	db.Model(&model.Category{}).First(&cg).Count(&count)
	if count == 0 {
		return nil, errors.New("标签不存在")
	}
	//查看名字是否被占用
	db.Model(&model.Category{}).Where("name=?", c.Name).Count(&count)
	if count != 0 {
		return nil, errors.New("标签已被占用")
	}
	cg.Name = c.Name
	db.Save(cg)
	if db.Error != nil {
		return nil, db.Error
	}
	return &cg, nil
}

/**
软删除标签
*/
func CategoryDelete(ids ...int) error {
	db := model.DB()
	db.Model(&model.Category{}).Where(ids).Update("soft_delete", 1)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

/**
恢复标签
*/
func CategoryRecover(ids ...int) error {
	db := model.DB()
	db.Model(&model.Category{}).Where(ids).Update("soft_delete", 0)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

/**
 * 获取所有的标签
 */
func CategoryGet() ([]*model.Category, error) {
	db := model.DB()
	cs := make([]*model.Category, 0)
	db.Order("created_at DESC").Find(&cs)
	if db.Error != nil {
		return nil, db.Error
	}
	return cs, nil
}
