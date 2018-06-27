package test

import (
	"math/rand"
	"media_framwork/model"
	"media_framwork/service"
	"strconv"
	"testing"
	"time"
)

//func TestModel(t *testing.T) {
//	model.Init()
//	t.Run("TestMUserCreate", TestMUserCreate)
//}

func TestMUserCreate(t *testing.T) {
	model.Init()
	rand.Seed(time.Now().Unix())
	s := strconv.Itoa(rand.Int())
	user := &model.User{
		Username: s,
		Password: "123456",
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

func TestMediaCreate(t *testing.T) {
	model.Init()
	m := &model.Media{
		Title:        "zzzz",
		Introduction: "ewqewqe",
	}
	//s := &model.Sharpness{Model : model.Model{ID : 1}}
	service.MediaCreate(m, []*model.MediaSharpness{
		&model.MediaSharpness{SharpnessId: 1, Uri: "zzzz"},
		&model.MediaSharpness{SharpnessId: 2, Uri: "zzzz"},
	}, []*model.Category{
		&model.Category{Model: model.Model{ID: 1}},
	})

	//model.DB().Where("user_id=?", user.ID).Delete(&model.UserRecord{})
}

func TestMediaUpdate(t *testing.T) {
	model.Init()

	//s := &model.Sharpness{Model : model.Model{ID : 1}}
	if _, err := service.MediaUpdate(5, []*model.MediaSharpness{
		&model.MediaSharpness{ID: 13, MediaID: 4, SharpnessId: 1, Uri: "qweqwe"},
		&model.MediaSharpness{ID: 14, MediaID: 6, SharpnessId: 2, Uri: "eqweqw"},
	}, []*model.Category{
		&model.Category{Model: model.Model{ID: 2}},
	}); err != nil {
		t.Error(err)
	}

	//model.DB().Where("user_id=?", user.ID).Delete(&model.UserRecord{})
}

func TestMUserLogin(t *testing.T) {
	model.Init()
	m := &model.User{
		Username: "dawndevil",
		Password: "123456",
	}
	if u, err := service.UserLogin(m); err != nil {
		t.Error(err)
	} else {
		t.Log(u)
	}

	//model.DB().Where("user_id=?", user.ID).Delete(&model.UserRecord{})
}
