package controller

import (
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"media_framwork/tool"
	"net/http"
	"strings"
)

/**
管理分类页面
*/
func PageCategoryManage(c *gin.Context) {
	var cs []model.Category
	curDB := model.DB()
	trash := c.DefaultQuery("trash", "0")
	//检索名字
	if name := c.DefaultQuery("name", ""); name != "" {
		curDB = curDB.Where("name=?", name)
	}
	p := model.DefaultPage(c)
	p.Find(&cs, curDB, trash != "0")
	c.HTML(http.StatusOK, "admin/categoryManage", h(gin.H{
		"data":  cs,
		"page":  p,
		"trash": trash,
	}, c))
}

/**
添加分类页面
*/
func PageCategoryAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/categoryAdd", h(gin.H{}, c))
}

/**
更新分类页面
*/
func PageCategoryUpdate(c *gin.Context) {
	id := c.Param("id")
	ct := model.Category{}
	if _, err := model.FindById(&ct, tool.GetInt(id)); err != nil {
		redirectError(c, "/404", err.Error())
		return
	}
	c.HTML(http.StatusOK, "admin/categoryUpdate", h(gin.H{
		"category": ct,
	}, c))
}

/**
分类创建控制器
*/
func CategoryCreate(c *gin.Context) {
	name := c.PostForm("name")
	if strings.TrimSpace(name) == "" {
		redirectError(c, "/admin/new_category", "标签名不能为空")
		return
	}
	ct := &model.Category{Name: name}
	if err := ct.Create(); err != nil {
		redirectError(c, "/admin/new_category", err.Error())
		return
	}
	redirectOK(c, "/admin/new_category", "创建标签成功")
}

/**
分类更新控制器
*/
func CategoryUpdate(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	ct := &model.Category{}
	ct.ID = tool.GetInt(id)
	ct.Name = name
	if err := ct.Update(); err != nil {
		redirectError(c, "/admin/category/update/"+id, err.Error())
		return
	}
	redirect(c, "/admin/category/update/"+id, nil)
}

/**
分类删除控制器
*/
func CategoryDelete(c *gin.Context) {
	ctStrings := c.PostFormArray("ct_ids")
	ctIds := tool.GetInts(ctStrings)
	if err := model.Delete(&model.Category{}, ctIds...); err != nil {
		redirectError(c, "/admin/manage_category", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_category", "删除成功")
}

/**
分类恢复控制器
*/
func CategoryRecover(c *gin.Context) {
	ctStrings := c.PostFormArray("ct_ids")
	ctIds := tool.GetInts(ctStrings)
	if err := model.Recover(&model.Category{}, ctIds...); err != nil {
		redirectError(c, "/admin/manage_category", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_category", "恢复成功")
}
