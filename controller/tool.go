package controller

import "github.com/gin-gonic/gin"

func result(code int, data interface{}, msg string) gin.H {
	return gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	}
}
