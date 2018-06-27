package controller

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"media_framwork/service"
	"net/http"
	"time"
)

/**
注册用户控制器
参数:
{
	"username" : ?
	"password" : ?,
	"passwordAgain" : ?
}
*/
func UserCreate(c *gin.Context) {
	var userDto struct {
		Username      string `json:"username" binding:"required"`
		Password      string `json:"password" binding:"required"`
		PasswordAgain string `json:"passwordAgain" binding:"required"`
		Sex           int    `json:"sex"`
	}

	if err := c.BindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
	} else {
		//验证两次的密码是否相同
		if userDto.Password != userDto.PasswordAgain {
			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, "两次密码不一致"))
			return
		}
		//根据Dto创建User
		user := &model.User{
			Username: userDto.Username,
			Password: hex.EncodeToString(md5.New().Sum([]byte(userDto.Password))),
			Sex:      userDto.Sex,
		}
		if u, err := service.UserCreate(user); err != nil {
			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
		} else {
			dto := UserReturnDto(u.CreatedAt, u.UpdatedAt, u.Username, u.Sex)
			c.JSON(http.StatusOK, result(http.StatusOK, dto, ""))
		}
	}
}

/**
用户登录控制器
参数:
{
	"username" : ?
	"password" : ?
}
*/
func UserLogin(c *gin.Context) {
	var userDto struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
	} else {
		user := &model.User{
			Username: userDto.Username,
			Password: userDto.Password,
		}
		if u, err := service.UserLogin(user); err != nil {
			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
		} else {
			dto := UserReturnDto(u.CreatedAt, u.UpdatedAt, u.Username, u.Sex)
			//设置session
			session := sessions.Default(c)
			session.Options(sessions.Options{
				MaxAge: int(24 * time.Hour),
				Path:   "/",
			})
			session.Set("userId", user.Username)
			session.Set("secret", md5EncodeWithSalt(u.Username, "dawndevil"))
			session.Save()

			c.JSON(http.StatusOK, result(http.StatusOK, dto, ""))
		}
	}
}

