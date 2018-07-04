package controller

import (
	"github.com/gin-gonic/gin"
	"dawn_media/conf"
	"dawn_media/model"
	"dawn_media/tool"
	"net/http"
	"strconv"
	"time"
	"errors"
	"strings"
)

/**
用户前台管理评论
 */
func PageODUserComment(c *gin.Context) {
	var authUser *model.User
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
	}
	var (
		page     = model.DefaultPage(c)
		curTime  = strconv.Itoa(int(time.Now().Unix()))
		comments = make([]*model.Comment, 0)
	)

	page.Find(
		&comments,
		model.CommentConditionDB(0, authUser.ID, c.DefaultQuery("curTime", curTime)),
		false)

	c.HTML(http.StatusOK, "ordinary/user_comment", h(gin.H{
		"comments": comments,
		"page":     page,
	}, c))
}
/**
用户后台管理评论
 */
func PageUserComment(c *gin.Context) {
	var (
		page     = model.DefaultPage(c)
		curTime  = strconv.Itoa(int(time.Now().Unix()))
		comments = make([]*model.Comment, 0)
	)

	page.Find(
		&comments,
		model.CommentConditionDB(0, 0, c.DefaultQuery("curTime", curTime)),
		false)
	c.HTML(http.StatusOK, "admin/manage_comment", h(gin.H{
		"comments": comments,
		"page":     page,
	}, c))
}

/**
管理员删除评论
 */
func CommentDelete(c *gin.Context) {
	commentDelete(c, "/admin/manage_comment")
}

/**
用户删除评论
 */
func ODCommentDelete(c *gin.Context) {
	commentDelete(c, "/ordinary/user_comment")
}

func commentDelete(c *gin.Context, uri string) {
	ids := tool.GetInts(c.PostFormArray("comment_ids"))
	if err := model.Delete(&model.Comment{}, ids...); err != nil {
		redirectError(c, uri, err.Error())
		return
	}
	redirectOK(c, uri, "删除成功")
}

/**
用户发表评论
 */
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
							"comment":   comment,
							"avatarMap": conf.C().AvatarMap,
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
		page    = model.DefaultPage(c)
		curTime = c.Query("curTime")
		media   = &model.Media{}
	)
	media.ID = tool.GetInt(c.DefaultQuery("id", "0"))
	if media.ID == 0 || !model.Exist(media, media.ID) {
		c.JSON(http.StatusBadRequest, j(http.StatusBadRequest, nil, "媒体不存在"))
		return
	}
	comments := make([]*model.Comment, 0)
	page.Find(&comments, model.CommentConditionDB(media.ID, 0, curTime), false)
	c.JSON(http.StatusOK, j(http.StatusOK, gin.H{
		"page":      page,
		"comments":  comments,
		"avatarMap": conf.C().AvatarMap,
	}, ""))
}
