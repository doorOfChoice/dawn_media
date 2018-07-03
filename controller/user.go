package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"media_framwork/model"
	"media_framwork/tool"
	"media_framwork/conf"
	"strconv"
	"path/filepath"
	"errors"
	"github.com/google/uuid"
	"github.com/gin-contrib/sessions"
)

/**
注册页面
 */
func PageRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "common/register", h(gin.H{}, c))
}

/**
登录页面
 */
func PageLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "common/login", h(gin.H{}, c))
}

/**
Admin管理页面
 */
func PageUserManage(c *gin.Context) {
	curDB := model.DB()
	trash := c.DefaultQuery("trash", "0")
	users := make([]*model.User, 0)
	page := model.DefaultPage(c)
	page.Find(&users, curDB, trash != "0")
	c.HTML(http.StatusOK, "admin/userManage", h(gin.H{
		"users": users,
		"trash": trash,
		"page" : page,
	}, c))
}

/**
Admin更新页面
 */
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

/**
Admin添加用户
 */
func PageUserAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/user_add", h(gin.H{}, c))
}

/**
Admin修改密码
 */
func PageODUserChangePwd(c *gin.Context) {
	c.HTML(http.StatusOK, "ordinary/update_pwd", h(gin.H{}, c))
}

/**
Admin修改Info
 */
func PageODUserChangeInfo(c *gin.Context) {
	c.HTML(http.StatusOK, "ordinary/update_info", h(gin.H{}, c))
}

/**
用户修改信息控制器
 */
func ODUserChangeInfo(c *gin.Context) {
	var authUser *model.User
	uri := "/ordinary/info_update"
	if t, ok := c.Get("authUser"); ok {
		authUser = t.(*model.User)
	}
	authUser.Nickname = c.PostForm("nickname")
	authUser.Password = ""

	v := MyValidator{}
	v.ValidateODUserUpdate(authUser)
	if redirectNotPass(c, uri, v) {
		return
	}
	if err := dealAvatar(c, authUser); err != nil {
		redirectError(c, uri, err.Error())
		return
	}
	if err := authUser.Update(); err != nil {
		redirectError(c, uri, err.Error())
		return
	}
	redirectOK(c, uri, "更新基本信息成功")
}

/**
用户修改密码控制器
 */
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

/**
用户添加控制器
 */
func UserAdd(c *gin.Context) {
	user := &model.User{}
	user.Username = c.PostForm("username")
	user.Nickname = c.PostForm("nickname")
	user.Password = c.PostForm("password")
	user.Authority = tool.GetInt(c.PostForm("authority"))
	user.Avatar = "default.jpg"
	//验证用户用户名和密码的合法性
	v := MyValidator{}
	v.ValidateUser(user)
	if redirectNotPass(c, "/admin/new_user", v) {
		return
	}

	var err error
	if err = dealAvatar(c, user); err == nil {
		if err = user.Create(); err == nil {
			redirectOK(c, "/admin/user/update/"+strconv.Itoa(user.ID), "创建用户成功")
			return
		}
	}
	redirectError(c, "/admin/new_user", err.Error())
}

/**
用户更新控制器
 */
func UserUpdate(c *gin.Context) {
	uri := "/admin/user/update/" + c.Param("id")
	id := tool.GetInt(c.Param("id"))
	user := &model.User{}
	user.ID = id
	user.Nickname = c.PostForm("nickname")
	user.Password = c.PostForm("password")
	user.Authority = tool.GetInt(c.PostForm("authority"))

	//验证用户用户名和密码的合法性
	v := MyValidator{}
	v.ValidateUserUpdate(user)
	if redirectNotPass(c, uri, v) {
		return
	}

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

/**
用户冻结
 */
func UserDelete(c *gin.Context) {
	userIds := tool.GetInts(c.PostFormArray("user_ids"))
	if err := model.Delete(&model.User{}, userIds...); err != nil {
		redirectError(c, "/admin/manage_user", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_user", "冻结成功")
}

/**
用户恢复
 */
func UserRecover(c *gin.Context) {
	userIds := tool.GetInts(c.PostFormArray("user_ids"))
	if err := model.Recover(&model.User{}, userIds...); err != nil {
		redirectError(c, "/admin/manage_user", err.Error())
		return
	}
	redirectOK(c, "/admin/manage_user", "恢复成功")
}

/**
注册控制器
 */
func Register(c *gin.Context) {
	uri := "/register"
	user := &model.User{}
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.Nickname = user.Username
	user.Avatar = "default.jpg"
	passwordAgain := c.PostForm("passwordAgain")

	//验证用户用户名和密码的合法性
	v := MyValidator{}
	v.Size(passwordAgain, 4, 20, "第二次输入的密码")
	v.ValidateUser(user)
	if redirectNotPass(c, uri, v) {
		return
	}

	if user.Password != passwordAgain {
		redirectError(c, uri, "两次密码输入不一致")
		return
	}
	if err := user.Create(); err != nil {
		redirectError(c, uri, err.Error())
		return
	}
	redirect(c, "/login", nil)
}

/**
登录控制器
 */
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
		//if u.Authority == model.ADMIN {
		//	redirectOK(c, "/admin", "登录成功")
		//}else {
		//	redirectOK(c, "/ordinary", "登录成功")
		//}
		redirectOK(c, "/", "登录成功")
	}
}

/**
注销
 */
func LogOut(c *gin.Context) {
	c.SetCookie(conf.C().SessionName, "", 0, "/", "", false, false)
	redirectOK(c, "/login", "")
}
