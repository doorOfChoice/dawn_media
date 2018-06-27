package controller

import (
	"github.com/gin-gonic/gin"
	"media_framwork/model"
	"media_framwork/service"
	"net/http"
)

/**
新建标签控制器
	参数
	{
		"name" : ?
	}
*/
func CategoryCreate(c *gin.Context) {
	var cgDto struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&cgDto); err != nil {
		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
	} else {
		cg := &model.Category{
			Name: cgDto.Name,
		}
		if cg2, err := service.CategoryCreate(cg); err != nil {
			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
		} else {
			c.JSON(http.StatusOK, result(http.StatusOK, cg2, ""))
		}
	}
}

/**
更新标签控制器
	参数
	host/v1/category/:id
	{
		"name" : ?
	}
*/
func CategoryBaseUpdate(c *gin.Context) {
	var cgDto struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&cgDto); err != nil {
		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
	} else {
		id := getInt(c.Query("id"))
		cg := &model.Category{
			Name: cgDto.Name,
		}
		if cg2, err := service.CategoryBaseUpdate(id, cg); err != nil {
			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
		} else {
			c.JSON(http.StatusOK, result(http.StatusOK, cg2, ""))
		}

	}
}

/**
软删除标签控制器
参数
host/v1/category?id=x&id=x&id=x

*/
func CategoryDelete(c *gin.Context) {
	var cgDto struct {
		Ids []int `form:"id" json:"id"`
	}
	if err := c.Bind(&cgDto); err != nil {
		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
	} else {
		if err := service.CategoryDelete(cgDto.Ids...); err != nil {
			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
		} else {
			c.JSON(http.StatusOK, result(http.StatusOK, nil, ""))
		}
	}
}

/**
恢复删除标签控制器
参数
host/v1/category?id=x&id=x&id=x
*/
func CategoryRecover(c *gin.Context) {
	var cgDto struct {
		Ids []int `form:"id" json:"id"`
	}
	if err := c.Bind(&cgDto); err != nil {
		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
	} else {
		if err := service.CategoryRecover(cgDto.Ids...); err != nil {
			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
		} else {
			c.JSON(http.StatusOK, result(http.StatusOK, nil, ""))
		}
	}
}

/**
获取所有分类
*/
func CategoryGet(c *gin.Context) {
	if cs, err := service.CategoryGet(); err != nil {
		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
	} else {
		c.JSON(http.StatusOK, result(http.StatusOK, cs, ""))
	}
}
