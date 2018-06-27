package main

import (
	"github.com/gin-gonic/gin"
	"media_framwork/controller"
	"media_framwork/model"
)

func main() {
	model.Init()
	router := gin.Default()

	router.POST("/user", controller.UserCreate)

	router.Run(":8080")
}
