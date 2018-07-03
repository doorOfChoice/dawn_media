package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"media_framwork/model"
)

/**
跳转404页面
这里的404指错误页面
 */
func Page404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "common/404", h(gin.H{}, c))
}

/**
管理员页面主页
 */
func PageAdminIndex(c *gin.Context) {
	mediaCount := model.Count(&model.Media{})
	categoryCount := model.Count(&model.Category{})
	userCount := model.Count(&model.User{})
	commentCount := model.Count(&model.Comment{})
	newUsers := model.GetAdminIndexNewUser()
	c.HTML(http.StatusOK, "admin/index", h(gin.H {
		"mediaCount" : mediaCount,
		"categoryCount" : categoryCount,
		"userCount" : userCount,
		"commentCount" : commentCount,
		"newUsers" : newUsers,
	},c ))
}