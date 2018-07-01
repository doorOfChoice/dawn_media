package controller

import (
	"github.com/gin-gonic/gin"
	"media_framwork/conf"
	"media_framwork/model"
	"media_framwork/tool"
	"net/http"
)

func PageODUserChangePwd(c *gin.Context) {
	c.HTML(http.StatusOK, "ordinary/update_pwd", h(gin.H{}, c))
}

func PageODUserChangeInfo(c *gin.Context) {
	c.HTML(http.StatusOK, "ordinary/update_info", h(gin.H{}, c))
}

func ODUserChangeInfo(c *gin.Context) {
	var authUser *model.User
	uri := "/ordinary/info_update"
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
	}
	authUser.Nickname = c.PostForm("nickname")
	v := MyValidator{}
	v.ValidateODUserUpdate(authUser)
	if redirectNotPass(c, uri, v) {
		return
	}
	if err := dealAvatar(c, authUser); err != nil {
		redirectError(c, uri, err.Error())
		return
	}
	redirectOK(c, uri, "更新基本信息成功")
}

func ODUserChangePwd(c *gin.Context) {
	var authUser *model.User
	uri := "/ordinary/pwd_update"
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
	}
	pwd := c.PostForm("origin")
	newPwd := c.PostForm("new_password")
	againPwd := c.PostForm("again_password")

	v := MyValidator{}
	v.ValidateODPwdUpdate(pwd, newPwd, againPwd)
	if redirectNotPass(c, uri, v) {
		return
	}

	if newPwd != againPwd {
		redirectError(c, uri, "两次输入的密码不一致")
		return
	}

	if tool.Md5EncodeWithSalt(pwd, conf.C().PassSalt) != authUser.Password {
		redirectError(c, uri, "旧密码错误")
		return
	}
	user := &model.User{}
	user.ID = authUser.ID
	user.Password = newPwd
	if err := user.Update(); err != nil {
		redirectError(c, uri, err.Error())
		return
	}
	redirectOK(c, uri, "更新密码成功")
}
