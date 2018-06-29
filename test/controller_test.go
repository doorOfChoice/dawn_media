package test

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"media_framwork/controller"
	"media_framwork/model"
	"testing"
)

func TestController(t *testing.T) {
	model.Init()
	store := cookie.NewStore([]byte("secret"))
	router := gin.Default()
	router.Use(sessions.Sessions("media_web", store))

	router.POST("/v1/user", controller.UserCreate)
	router.POST("/v1/user/session", controller.UserLogin)

	router.POST("/v1/category", controller.CategoryCreate)
	router.DELETE("/v1/category", controller.CategoryDelete)
	router.DELETE("/v1/trash/category", controller.CategoryRecover)

	//router.POST("/v1/media", controller.MediaCreate)
	router.Run(":8080")

}
