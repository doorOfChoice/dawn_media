package test

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"dawn_media/controller"
	"dawn_media/model"
	"testing"
)

func TestController(t *testing.T) {
	model.Init()
	store := cookie.NewStore([]byte("secret"))
	router := gin.Default()
	router.Use(sessions.Sessions("media_web", store))

	router.POST("/v1/category", controller.CategoryCreate)
	router.DELETE("/v1/category", controller.CategoryDelete)
	router.DELETE("/v1/trash/category", controller.CategoryRecover)

	//router.POST("/v1/media", controller.MediaCreate)
	router.Run(":8080")

}

func TestValidate(t *testing.T) {
	name := "1a"

	v := validation.Validation{}
	v.MinSize(name, 2, "name").Message("最小长度应该大于2")

	if v.HasErrors() {
		for _, err := range v.Errors {
			t.Log(err)
		}
	}
}
