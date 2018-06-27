package test

import (
	"testing"
	"media_framwork/model"
	"github.com/gin-gonic/gin"
	"media_framwork/controller"
)

func TestCUserCreate(t *testing.T) {
	model.Init()

	router := gin.Default()

	router.POST("/user", controller.UserCreate)

	router.Run(":8080")
}