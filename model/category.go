package model

import "errors"

/**
创建分类
*/
func (c *Category) Create() error {
	count := 0
	db.Model(Category{}).Where("name=?", c.Name).Count(&count)
	if count != 0 {
		return errors.New("分类已经存在")
	}
	db.Save(c)
	return db.Error
}

/**
更新分类
*/
func (c *Category) Update() error {
	if u, err := FindById(&Category{}, c.ID); err != nil {
		return err
	} else {
		ct := u.(*Category)
		count := 0
		db.Model(c).Where("name=?", c.Name).Count(&count)
		if count != 0 {
			return errors.New("标签名已被占用")
		}
		ct.Name = c.Name
		db.Save(ct)
		return db.Error
	}
}

/**
获取没有删除的所有分类
 */
func GetCategories() []*Category {
	categories := make([]*Category, 0)
	db.Where("soft_delete=?", 1).Find(&categories)
	return categories
}
