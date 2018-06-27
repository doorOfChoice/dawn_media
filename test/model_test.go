package test

import (
	"media_framwork/model"
	"media_framwork/service"
	"testing"
	"math/rand"
	"time"
	"strconv"
)

func TestModel(t *testing.T) {
	model.Init()
	t.Run("TestMUserCreate", TestMUserCreate)
}

func TestMUserCreate(t *testing.T) {
	rand.Seed(time.Now().Unix())
	s:= strconv.Itoa(rand.Int())
	user := &model.User{
		Username: s,
		Password: "12321",
	}
	if u, err := service.UserCreate(user); err != nil {
		t.Error(err)
	} else {
		t.Log("第一次创建: ", u)
	}

	if u, err := service.UserCreate(user); err != nil {
		t.Error(err)
	} else {
		t.Log(u)
	}
	//model.DB().Where("user_id=?", user.ID).Delete(&model.UserRecord{})
}
