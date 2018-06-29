package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"media_framwork/conf"
	"media_framwork/model"
	"media_framwork/tool"
	"net/http"
	"path/filepath"
)

func PageMediaAdd(c *gin.Context) {
	categories := make([]*model.Category, 0)
	model.DB().Select("id,name").Find(&categories)
	c.HTML(http.StatusOK, "admin/media_add", h(gin.H{
		"categories": categories,
	}, c))
}
func PageMediaUpdate(c *gin.Context) {
	id := tool.GetInt(c.Param("id"))
	media := &model.Media{}
	media.ID = id
	if err := media.Get(); err != nil {
		redirectError(c, "/404", err.Error())
		return
	}
	categories := make([]*model.Category, 0)
	model.DB().Select("id,name").Find(&categories)
	c.HTML(http.StatusOK, "admin/media_update", h(gin.H{
		"categories": categories,
		"media":      media,
	}, c))
}
func MediaUpdate(c *gin.Context) {
	media := &model.Media{}
	media.ID = tool.GetInt(c.Param("id"))
	media.Title = c.PostForm("title")
	media.Introduction = c.PostForm("introduction")
	ctIds := tool.GetInts(c.PostFormArray("categories"))
	for _, id := range ctIds {
		ct := &model.Category{}
		ct.ID = id
		media.Categories = append(media.Categories, ct)
	}
	var s string
	var err error
	if s, err = dealMediaFile(c, "file_s360"); err == nil {
		media.S360 = s
		if s, err = dealMediaFile(c, "file_s480"); err == nil {
			media.S480 = s
			if s, err = dealMediaFile(c, "file_s720"); err == nil {
				media.S720 = s
				if s, err = dealMediaFile(c, "file_s1080"); err == nil {
					media.S1080 = s
					if err = media.Update(); err == nil {
						redirectOK(c, "/admin/media/update/"+c.Param("id"), "更新媒体成功")
						return
					}
				}
			}
		}
	}
	redirectError(c, "/admin/update/media"+c.Param("id"), err.Error())
}
func MediaAdd(c *gin.Context) {
	media := &model.Media{}
	media.Title = c.PostForm("title")
	media.Introduction = c.PostForm("introduction")
	ctIds := tool.GetInts(c.PostFormArray("categories"))
	for _, id := range ctIds {
		ct := &model.Category{}
		ct.ID = id
		media.Categories = append(media.Categories, ct)
	}
	var s string
	var err error
	if s, err = dealMediaFile(c, "file_s360"); err == nil {
		media.S360 = s
		if s, err = dealMediaFile(c, "file_s480"); err == nil {
			media.S480 = s
			if s, err = dealMediaFile(c, "file_s720"); err == nil {
				media.S720 = s
				if s, err = dealMediaFile(c, "file_s1080"); err == nil {
					media.S1080 = s
					if err = media.Create(); err == nil {
						redirectOK(c, "/admin/new_media", "创建媒体成功")
						return
					}
				}
			}
		}
	}
	redirectError(c, "/admin/new_media", err.Error())
}

func dealMediaFile(c *gin.Context, formName string) (string, error) {
	parent := conf.C().UploadDir
	if info, err := c.FormFile(formName); err != nil {
		return "", nil
	} else {
		if info.Size == 0 {
			return "", nil
		}
		ext := filepath.Ext(info.Filename)
		if ext != ".mp4" && ext != ".avi" && ext != ".ts" {
			return "", errors.New("文件格式不支持")
		}

		filename := uuid.New().String() + ext
		if err := c.SaveUploadedFile(info, parent+filename); err != nil {
			return "", err
		}
		return filename, nil
	}

}
