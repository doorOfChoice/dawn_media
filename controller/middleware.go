package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"media_framwork/model"
)

/**
管理员后台管理中间件
*/
func MiddlewareAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if id, ok := session.Get("userId").(int); ok {
			user := &model.User{}
			user.ID = id
			if err := user.Get(); err == nil {
				if user.Authority != model.ADMIN {
					redirectError(c, "/404", "用户无权限访问")
					return
				}
				if user.SoftDelete == 2 {
					redirectError(c, "/404", "用户被冻结")
					return
				}
				c.Set("authUser", user)
				c.Next()
				return
			}
		}
		redirectError(c, "/404", "权限验证失败")
	}
}

/**
普通用户后台管理中间件
*/
func MiddlewareOrdinaryAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if id, ok := session.Get("userId").(int); ok {
			user := &model.User{}
			user.ID = id
			if err := user.Get(); err == nil {
				if user.SoftDelete == 2 {
					redirectError(c, "/404", "用户被冻结")
					return
				}
				c.Set("authUser", user)
				c.Next()
				return
			}
		}
		redirectError(c, "/404", "权限验证失败")
	}
}

/**
前台中间件
*/
func MiddlewareIndexAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if id, ok := session.Get("userId").(int); ok {
			user := &model.User{}
			user.ID = id
			if err := user.Get(); err == nil {
				if user.SoftDelete == 2 {
					redirectError(c, "/404", "用户被冻结")
					return
				}
				c.Set("authUser", user)
			}
		}
		c.Next()
	}
}
