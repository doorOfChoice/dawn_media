package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"dawn_media/conf"
	"dawn_media/model"
	"dawn_media/tool"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"log"
)

type mediaAttributeDto struct {
	id          int
	file        *multipart.FileHeader
	description string
}

/**
添加媒体页面
*/
func PageMediaAdd(c *gin.Context) {
	categories := make([]*model.Category, 0)
	model.DB().Select("id,name").Where("soft_delete=1").Find(&categories)
	c.HTML(http.StatusOK, "admin/media_add", h(gin.H{
		"categories": categories,
	}, c))
}

/**
更新媒体页面
*/
func PageMediaUpdate(c *gin.Context) {
	id := tool.GetInt(c.Param("id"))
	media := &model.Media{}
	media.ID = id
	if err := media.Get(); err != nil {
		redirectError(c, "/404", err.Error())
		return
	}

	categories := make([]*model.Category, 0)
	model.DB().Select("id,name").Where("soft_delete=1").Find(&categories)
	c.HTML(http.StatusOK, "admin/media_update", h(gin.H{
		"categories": categories,
		"media":      media,
	}, c))
}

/**
管理媒体控制器
*/
func PageMediaManage(c *gin.Context) {
	medias := make([]*model.Media, 0)
	page := model.DefaultPage(c)
	curDB := model.DB()
	trash := c.DefaultQuery("trash", "0")
	if title := c.DefaultQuery("title", ""); title != "" {
		curDB = curDB.Where("title=?", title)
	}
	if id := c.DefaultQuery("id", ""); id != "" {
		curDB = curDB.Where("id=?", id)
	}
	page.Find(&medias, curDB, trash != "0")
	c.HTML(http.StatusOK, "admin/mediaManage", h(gin.H{
		"data":  medias,
		"trash": trash,
		"page":  page,
	}, c))
}

/**
用户点过赞视频
 */
func PageODStaredMedias(c *gin.Context) {
	var authUser *model.User
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
	}
	page := model.DefaultPage(c)
	medias := make([]*model.Media, 0)
	page.Find(&medias, model.StaredMediasDB(authUser.ID))
	c.HTML(http.StatusOK, "ordinary/stared_medias", h(gin.H{
		"medias" : medias,
		"page" : page,
	}, c))
}

/**
更新媒体控制器
*/
func MediaUpdate(c *gin.Context) {
	uri := "/admin/media/update/" + c.Param("id")
	media := &model.Media{}
	media.ID = tool.GetInt(c.Param("id"))
	media.Title = c.PostForm("title")
	media.Introduction = c.PostForm("introduction")

	v := MyValidator{}
	v.ValidateMedia(media)
	if redirectNotPass(c, uri, v) {
		return
	}

	ctIds := tool.GetInts(c.PostFormArray("ct_ids"))
	for _, id := range ctIds {
		ct := &model.Category{}
		ct.ID = id
		media.Categories = append(media.Categories, ct)
	}

	var err error
	//跳转网址
	if err = dealMediaFile(c, media, hookFile(c)); err == nil {
		if err = media.Update(); err == nil {
			redirectOK(c, uri, "更新成功")
			return
		}
	}
	redirectError(c, uri, err.Error())
}

/**
删除媒体控制器
*/
func MediaDelete(c *gin.Context) {
	mediaIds := tool.GetInts(c.PostFormArray("media_ids"))
	if err := model.Delete(&model.Media{}, mediaIds...); err != nil {
		redirectError(c, "/admin/manage_media", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_media", "删除媒体成功")
}

/**
恢复媒体控制器
*/
func MediaRecover(c *gin.Context) {
	mediaIds := tool.GetInts(c.PostFormArray("media_ids"))
	if err := model.Recover(&model.Media{}, mediaIds...); err != nil {
		redirectError(c, "/admin/manage_media", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_media", "恢复媒体成功")
}

/**
添加媒体控制器
*/
func MediaAdd(c *gin.Context) {
	media := &model.Media{}
	media.Title = c.PostForm("title")
	media.Introduction = c.PostForm("introduction")

	v := MyValidator{}
	v.ValidateMedia(media)
	if redirectNotPass(c, "/admin/new_media", v) {
		return
	}

	ctIds := tool.GetInts(c.PostFormArray("ct_ids"))
	for _, id := range ctIds {
		ct := &model.Category{}
		ct.ID = id
		media.Categories = append(media.Categories, ct)
	}

	var err error
	if err = dealMediaFile(c, media, hookFile(c)); err == nil {
		if err = media.Create(); err == nil {
			redirectOK(c, "/admin/media/update/"+strconv.Itoa(media.ID), "创建媒体成功")
			return
		}
	}
	redirectError(c, "/admin/new_media", err.Error())
}

func hookFile(c *gin.Context) []mediaAttributeDto {
	ms := make([]mediaAttributeDto, 4)
	if u, err := c.FormFile("file_s360"); err == nil {
		ms[0].description = "360P"
		ms[0].file = u
		ms[0].id = tool.GetInt(c.DefaultPostForm("file_s360_id", "0"))
	}
	if u, err := c.FormFile("file_s480"); err == nil {
		ms[1].description = "480P"
		ms[1].file = u
		ms[1].id = tool.GetInt(c.DefaultPostForm("file_s480_id", "0"))
	}
	if u, err := c.FormFile("file_s720"); err == nil {
		ms[2].description = "720P"
		ms[2].file = u
		ms[2].id = tool.GetInt(c.DefaultPostForm("file_s720_id", "0"))
	}
	if u, err := c.FormFile("file_s1080"); err == nil {
		ms[3].description = "1080P"
		ms[3].file = u
		ms[3].id = tool.GetInt(c.DefaultPostForm("file_s1080_id", "0"))
	}

	return ms
}

/**
处理媒体文件
*/
func dealMediaFile(c *gin.Context, m *model.Media, infos []mediaAttributeDto) error {
	if cover, err := c.FormFile("cover"); err == nil && cover.Size != 0 {
		filename := uuid.New().String() + filepath.Ext(cover.Filename)
		if err := c.SaveUploadedFile(cover, conf.C().CoverDir+filename); err == nil {
			m.Cover = filename
		}
	}
	for _, info := range infos {
		//后缀
		ext := filepath.Ext(info.file.Filename)
		//前缀
		base := filepath.Base(info.file.Filename)
		//TODO 验证视频格式
		//单纯判断名字是否符合
		if strings.TrimSpace(base) != "." &&
			ext != ".mp4" &&
			ext != ".avi" &&
			ext != ".ts" {
			return errors.New(info.description + "文件格式不支持")
		}
	}
	for _, info := range infos {
		ma := &model.MediaAttribute{Description: info.description}
		log.Println(ma)
		m.MediaAttributes = append(m.MediaAttributes, ma)
		base := filepath.Base(info.file.Filename)
		if info.file == nil ||
			info.file.Size == 0 ||
			strings.TrimSpace(base) == "." {
			continue
		}
		ext := filepath.Ext(info.file.Filename)
		filename := uuid.New().String() + ext
		if err := c.SaveUploadedFile(info.file, conf.C().MediaDir+filename); err != nil {
			return err
		}
		ma.ID = info.id
		ma.Filename = info.file.Filename
		ma.Uri = filename
	}
	return nil
}
