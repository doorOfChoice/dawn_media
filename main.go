package main

import (
	"dawn_media/conf"
	"dawn_media/controller"
	"dawn_media/model"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"time"
	"log"
)

func format(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

func formatDetail(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func initDir() {
	//初始化目录
	if err := os.MkdirAll(conf.C().AvatarDir, 0755); err != nil {
		log.Println(err)
	}
	if err := os.MkdirAll(conf.C().MediaDir, 0755); err!= nil {
		log.Println(err)
	}
	if err := os.MkdirAll(conf.C().CoverDir, 0755); err != nil {
		log.Println(err)
	}
}

func main() {
	conf.Init()
	//初始化目录
	initDir()
	model.Init()
	store := cookie.NewStore([]byte(conf.C().PassSalt))
	r := gin.Default()
	r.Use(sessions.Sessions(conf.C().SessionName, store))
	//先设置函数，因为解析模板的时候回解析函数
	r.SetFuncMap(template.FuncMap{
		"formatDate":   format,
		"formatDetail": formatDetail,
	})
	r.LoadHTMLGlob("views/**/*")
	r.Static("/static", "./static")

	//Admin页面
	auth := r.Group("/admin")
	auth.Use(controller.MiddlewareAdminAuth())
	{
		auth.GET("/", controller.PageAdminIndex)
		auth.GET("/new_category", controller.PageCategoryAdd)
		auth.GET("/manage_category", controller.PageCategoryManage)
		auth.GET("/category/update/:id", controller.PageCategoryUpdate)
		auth.POST("/new_category", controller.CategoryCreate)
		auth.POST("/category/update/:id", controller.CategoryUpdate)
		auth.POST("/category/delete", controller.CategoryDelete)
		auth.POST("/category/recover", controller.CategoryRecover)

		auth.GET("/new_media", controller.PageMediaAdd)
		auth.GET("/media/update/:id", controller.PageMediaUpdate)
		auth.GET("/manage_media", controller.PageMediaManage)
		auth.POST("/media/update/:id", controller.MediaUpdate)
		auth.POST("/new_media", controller.MediaAdd)
		auth.POST("/media/delete", controller.MediaDelete)
		auth.POST("/media/recover", controller.MediaRecover)

		auth.GET("/new_user", controller.PageUserAdd)
		auth.GET("/manage_user", controller.PageUserManage)
		auth.GET("/user/update/:id", controller.PageUserUpdate)
		auth.POST("/new_user", controller.UserAdd)
		auth.POST("/user/update/:id", controller.UserUpdate)
		auth.POST("/user/delete", controller.UserDelete)
		auth.POST("/user/recover", controller.UserRecover)

		auth.GET("/manage_comment", controller.PageUserComment)
		auth.POST("/comment/delete", controller.CommentDelete)
	}

	ordinary := r.Group("/ordinary")
	ordinary.Use(controller.MiddlewareOrdinaryAuth())
	{
		ordinary.GET("/", controller.PageODUserChangeInfo)
		ordinary.GET("/info_update", controller.PageODUserChangeInfo)
		ordinary.GET("/pwd_update", controller.PageODUserChangePwd)
		ordinary.GET("/user_record", controller.PageODUserRecord)
		ordinary.GET("/user_comment", controller.PageODUserComment)
		ordinary.GET("/stared_medias", controller.PageODStaredMedias)
		ordinary.POST("/info_update", controller.ODUserChangeInfo)
		ordinary.POST("/pwd_update", controller.ODUserChangePwd)
		ordinary.POST("/comment/delete", controller.ODCommentDelete)
	}

	front := r.Group("/")
	front.Use(controller.MiddlewareIndexAuth())
	{
		front.GET("/", controller.PageFrontIndex)
		front.GET("/medias", controller.PageFrontMedias)
		front.GET("/single", controller.PageFrontSingle)
		front.GET("/comments", controller.GetComments)
		front.POST("/comments", controller.CommentCreate)
		front.POST("/media/star", controller.StarToggle)
	}

	//普通页面
	r.GET("/404", controller.Page404)
	r.GET("/register", controller.PageRegister)
	r.GET("/login", controller.PageLogin)
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.POST("/logout", controller.LogOut)
	r.Run(conf.C().Address)
}
