package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"media_framwork/model"
)

func Page404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "common/404", h(gin.H{}, c))
}

func PageAdminIndex(c *gin.Context) {
	mediaCount := model.Count(&model.Media{})
	categoryCount := model.Count(&model.Category{})
	userCount := model.Count(&model.User{})
	c.HTML(http.StatusOK, "admin/index", h(gin.H {
		"mediaCount" : mediaCount,
		"categoryCount" : categoryCount,
		"userCount" : userCount,
	},c ))
}