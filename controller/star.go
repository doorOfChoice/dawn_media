package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"net/http"
)

func StarToggle(c *gin.Context) {
	var (
		authUser *model.User
		err      = errors.New("请登录")
	)
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
		var star struct {
			MediaID int `json:"media_id,string"`
		}
		if e := c.ShouldBindJSON(&star); e == nil {
			if create, count, e := model.ToggleStar(authUser.ID, star.MediaID); e == nil {
				c.JSON(http.StatusOK, j(http.StatusOK, gin.H{
					"create": create,
					"count":  count,
				}, ""))
				return
			} else {
				err = e
			}
		} else {
			err = e
		}
	}
	c.JSON(http.StatusBadRequest, j(http.StatusBadRequest, nil, err.Error()))
}
