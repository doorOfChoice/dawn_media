package controller

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"net/http"
)

type MyValidator struct {
	validation.Validation
}

//v应该为指针，否则是向副本里面加入的参数，而不是向实体
func (v *MyValidator) Size(i interface{}, min int, max int, message ...string) {
	name := "参数"
	if len(message) == 1 {
		name = message[0]
	}
	v.MinSize(i, min, "key1").Message(fmt.Sprintf("%s应该大于%d", name, min))
	v.MaxSize(i, max, "key2").Message(fmt.Sprintf("%s应该小于%d", name, max))
}

func (v *MyValidator) ValidateMedia(media *model.Media) {
	v.Size(media.Title, 1, 40, "媒体标题")
	v.Size(media.Introduction, 1, 300, "媒体介绍")
}

func (v *MyValidator) ValidateUser(user *model.User) {
	v.AlphaDash(user.Username, "key1").Message("用户名只能包含数字字母-_")
	v.Size(user.Username, 4, 20, "用户名")
	v.Size(user.Password, 4, 30, "密码")
}

func (v *MyValidator) ValidateUserUpdate(user *model.User) {
	if user.Password != "" {
		v.Size(user.Password, 4, 30, "密码")
	}
	v.Size(user.Nickname, 1, 20, "用户昵称")
}
func (v *MyValidator) ValidateODUserUpdate(user *model.User) {
	v.Size(user.Nickname, 1, 20, "用户昵称")
}
func (v *MyValidator) ValidateODPwdUpdate(o, n, a string) {
	v.Size(o, 4, 30, "原始密码")
	v.Size(n, 4, 30, "新密码")
	v.Size(a, 4, 30, "重复的密码")
}
/**
在原始返回数据里封装错误和成功消息
*/
func h(gh gin.H, c *gin.Context) gin.H {
	gh["error"] = c.Query("error")
	gh["success"] = c.Query("success")
	if authUser, ok := c.Get("authUser"); ok {
		gh["authUser"] = authUser
	}
	return gh
}

/**
重定向
*/
func redirect(c *gin.Context, website string, h gin.H) {
	query := ""
	buf := bytes.NewBufferString(query)
	ok := false
	for k, v := range h {
		if ok {
			buf.WriteString("&")
		}
		ok = true
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v.(string))
	}
	c.Redirect(http.StatusMovedPermanently, website+"?"+buf.String())
}

func redirectOK(c *gin.Context, website string, success string) {
	redirect(c, website, gin.H{
		"success": success,
	})
}
func redirectError(c *gin.Context, website string, err string) {
	redirect(c, website, gin.H{
		"error": err,
	})
}
func redirectNotPass(c *gin.Context, website string, v MyValidator) bool {
	info := ""
	if v.HasErrors() {
		for _, err := range v.Errors {
			info += "<" + err.Error() + ">"
		}
		redirectError(c, website, info)
		return true
	}
	return false
}
