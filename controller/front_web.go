package controller

import (
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"media_framwork/tool"
	"net/http"
)

func PageFrontIndex(c *gin.Context) {
	var (
		randomMedias = model.GetIndexRandomMedia()
		hotMedias    = model.GetIndexHotMedia()
		newMedias    = model.GetIndexNewMedia()
		newComments = model.GetIndexNewComments()
		categories   = model.GetCategories()
	)
	c.HTML(http.StatusOK, "front/index", h(gin.H{
		"randomMedias": randomMedias,
		"hotMedias":    hotMedias,
		"newMedias":    newMedias,
		"categories":   categories,
		"newComments" : newComments,
	}, c))
}

func PageFrontMedias(c *gin.Context) {
	var (
		byTime     = tool.GetInt(c.DefaultQuery("byTime", "0"))
		byHot      = tool.GetInt(c.DefaultQuery("byHot", "0"))
		title      = c.DefaultQuery("title", "")
		category   = tool.GetInt(c.DefaultQuery("category", "0"))
		page       = model.DefaultPage(c)
		categories = model.GetCategories()
		curDB      = model.ConditionSelectMediaDB(byTime, byHot, category, title)
	)

	medias := make([]*model.Media, 0)
	page.Find(&medias, curDB)
	c.HTML(http.StatusOK, "front/medias", h(gin.H{
		"medias":     medias,
		"page":       page,
		"categories": categories,
	}, c))
}

func PageFrontSingle(c *gin.Context) {
	var (
		categories = model.GetCategories()
		newMedias  = model.GetIndexNewMedia()
		authUser *model.User
	)

	media := &model.Media{}
	media.ID = tool.GetInt(c.DefaultQuery("id", "0"))
	if err := media.Get(); err != nil {
		redirectError(c, "/404", "视频不存在")
		return
	}
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
		model.UserRecordUpdate(authUser.ID, media.ID)
	}
	//更新历史记录
	model.UserRecordUpdate(authUser.ID, media.ID)
	c.HTML(http.StatusOK, "front/single", h(gin.H{
		"categories": categories,
		"newMedias":  newMedias,
		"media":      media,
	}, c))
}
