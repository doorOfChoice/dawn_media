package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"net/http"
	"media_framwork/conf"
	"strings"
	"media_framwork/tool"
)

func CommentCreate(c *gin.Context) {
	var (
		authUser *model.User
		err      = errors.New("无权限发表评论")
	)
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
		comment := &model.Comment{}
		comment.UserID = authUser.ID
		if e := c.ShouldBindJSON(comment); e == nil {
			if len(strings.TrimSpace(comment.Content)) < 1 {
				err = errors.New("评论不能为空")
			} else {
				if e := comment.Create(); e == nil {
					if e := comment.Get(); e == nil {
						c.JSON(http.StatusOK, j(http.StatusOK, gin.H{
							"comment" : comment,
							"avatarMap" : conf.C().AvatarMap,
						}, ""))
						return
					} else {
						err = e
					}
				} else {
					err = e
				}
			}
		} else {
			err = e
		}

	}
	c.JSON(http.StatusBadRequest, j(http.StatusBadRequest, "", err.Error()))
}

/**
获取评论，根据时间戳和媒体id获取，时间戳为了防止重复加载
 */
func GetComments(c *gin.Context) {
	var (
		page  = model.DefaultPage(c)
		curTime = c.Query("curTime")
		media = &model.Media{}
	)
	media.ID = tool.GetInt(c.DefaultQuery("id", "0"))
	if media.ID == 0 || !model.Exist(media, media.ID) {
		c.JSON(http.StatusBadRequest, j(http.StatusBadRequest, nil, "媒体不存在"))
		return
	}
	comments := make([]*model.Comment, 0)
	page.Find(&comments, model.CommentConditionDB(media.ID, curTime), false)
	c.JSON(http.StatusOK, j(http.StatusOK, gin.H{
		"page":     page,
		"comments": comments,
		"avatarMap" : conf.C().AvatarMap,
	}, ""))
}
