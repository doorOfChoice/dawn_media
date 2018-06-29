package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"media_framwork/conf"
	"media_framwork/controller"
	"media_framwork/model"
	"net/http"
	"time"
)

func format(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

func main() {
	model.Init()
	conf.Init()
	r := gin.Default()
	//先设置函数，因为解析模板的时候回解析函数
	r.SetFuncMap(template.FuncMap{
		"formatDate": format,
	})
	r.LoadHTMLGlob("views/**/*")
	r.Static("/static", "./static")

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin/index", nil)
	})
	r.GET("/404", controller.Page404)
	r.GET("/admin/new_category", controller.PageCategoryAdd)
	r.GET("/admin/manage_category", controller.PageCategoryManage)
	r.GET("/admin/category/update/:id", controller.PageCategoryUpdate)
	r.POST("/admin/new_category", controller.CategoryCreate)
	r.POST("/admin/category/update/:id", controller.CategoryUpdate)
	r.POST("/admin/category/delete", controller.CategoryDelete)
	r.POST("/admin/category/recover", controller.CategoryRecover)

	r.GET("/admin/new_media", controller.PageMediaAdd)
	r.GET("/admin/media/update/:id", controller.PageMediaUpdate)
	r.POST("/admin/media/update/:id", controller.MediaUpdate)
	r.POST("/admin/new_media", controller.MediaAdd)

	r.Run(conf.C().Address)
}
