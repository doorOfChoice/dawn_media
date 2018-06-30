package controller

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"media_framwork/conf"
	"media_framwork/model"
	"media_framwork/tool"
	"net/http"
	"path/filepath"
	"strconv"
)

func PageRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "common/register", h(gin.H{}, c))
}

func PageLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "common/login", h(gin.H{}, c))
}

func PageUserManage(c *gin.Context) {
	curDB := model.DB()
	trash := c.DefaultQuery("trash", "0")
	users := make([]*model.User, 0)
	p := model.DefaultPage(c)
	p.Find(&users, curDB, trash != "0")
	for _, user := range users {
		user.Avatar = conf.C().AvatarMap + user.Avatar
	}
	c.HTML(http.StatusOK, "admin/userManage", h(gin.H{
		"users": users,
		"trash": trash,
	}, c))
}

func PageUserUpdate(c *gin.Context) {
	id := tool.GetInt(c.Param("id"))
	user := &model.User{}
	user.ID = id
	if err := user.Get(); err != nil {
		redirectError(c, "/404", err.Error())
		return
	}
	c.HTML(http.StatusOK, "admin/user_update", h(gin.H{
		"user": user,
	}, c))
}

func PageUserAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/user_add", h(gin.H{}, c))
}

func UserAdd(c *gin.Context) {
	user := &model.User{}
	user.Username = c.PostForm("username")
	user.Nickname = c.PostForm("nickname")
	user.Password = c.PostForm("password")
	user.Authority = tool.GetInt(c.PostForm("authority"))
	var err error
	if err = dealAvatar(c, user); err == nil {
		if err = user.Create(); err == nil {
			redirectOK(c, "/admin/user/update/"+strconv.Itoa(user.ID), "创建用户成功")
			return
		}
	}
	redirectError(c, "/admin/new_user", err.Error())
}

func UserUpdate(c *gin.Context) {
	id := tool.GetInt(c.Param("id"))
	user := &model.User{}
	user.ID = id
	user.Nickname = c.PostForm("nickname")
	user.Password = c.PostForm("password")
	user.Authority = tool.GetInt(c.PostForm("authority"))
	uri := "/admin/user/update/" + c.Param("id")
	var err error
	if err = dealAvatar(c, user); err == nil {
		if err = user.Update(); err == nil {
			redirectOK(c, uri, "更新用户成功")
			return
		}
	}
	redirectError(c, uri, err.Error())
}

func dealAvatar(c *gin.Context, user *model.User) error {
	if f, err := c.FormFile("avatar"); err != nil {
		return err
	} else {
		if f.Size == 0 {
			return nil
		}
		ext := filepath.Ext(f.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			return errors.New("图片格式不支持")
		}
		prefix := uuid.New().String()
		filename := prefix + ext
		if err := c.SaveUploadedFile(f, conf.C().AvatarDir+filename); err != nil {
			return err
		}
		user.Avatar = filename
		return nil
	}
}

func UserDelete(c *gin.Context) {
	userIds := tool.GetInts(c.PostFormArray("user_ids"))
	if err := model.Delete(&model.User{}, userIds...); err != nil {
		redirectError(c, "/admin/manage_user", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_user", "冻结成功")
}

func UserRecover(c *gin.Context) {
	userIds := tool.GetInts(c.PostFormArray("user_ids"))
	if err := model.Recover(&model.User{}, userIds...); err != nil {
		redirectError(c, "/admin/manage_user", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_user", "恢复成功")
}

func Register(c *gin.Context) {
	user := &model.User{}
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.Nickname = user.Username
	user.Avatar = "default.jpg"
	passwordAgain := c.PostForm("passwordAgain")
	if user.Password != passwordAgain {
		redirectError(c, "/register", "两次密码输入不一致")
		return
	}
	if err := user.Create(); err != nil {
		redirectError(c, "/register", err.Error())
		return
	}
	redirect(c, "/login", nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if u, err := model.FindUserByLogin(username, password); err != nil {
		redirectError(c, "/login", err.Error())
	} else {
		session := sessions.Default(c)
		session.Set("userId", u.ID)
		session.Save()
		if u.SoftDelete == 2 {
			redirectError(c, "/login", "用户已被冻结")
			return
		}
		if u.Authority == model.ADMIN {
			redirectOK(c, "/admin", "登录成功")
		}
	}
}

func LogOut(c *gin.Context) {
	c.SetCookie(conf.C().SessionName, "", 0, "/", "", false, false)
	redirectOK(c, "/login", "")
}
