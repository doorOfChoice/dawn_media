package controller

import (
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"net/http"
)

func PageODUserRecord(c *gin.Context) {
	var user *model.User
	if t, ok := c.Get("authUser"); ok {
		user = t.(*model.User)
	}
	curDB := model.UserRecordDB(user.ID)
	page := model.DefaultPage(c)
	records := make([]*model.UserRecord, 0)
	page.Find(&records, curDB)

	c.HTML(http.StatusOK, "ordinary/user_record", h(gin.H{
		"records": records,
		"page":    page,
	}, c))
}
