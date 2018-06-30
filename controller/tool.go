package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
返回固定的模式
*/
func result(code int, data interface{}, msg string) gin.H {
	return gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	}
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
