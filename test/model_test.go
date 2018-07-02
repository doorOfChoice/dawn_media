package test

import (
	"media_framwork/model"
	"testing"
	"github.com/google/uuid"
	"log"
	"media_framwork/conf"
	"path/filepath"
)

//func TestModel(t *testing.T) {
//	model.Init()
//	t.Run("TestMUserCreate", TestMUserCreate)
//}

func TestMUserCreate(t *testing.T) {
	model.Init()
	log.Println(model.GetRandomMedia())
	log.Println(model.GetHotMedia())
}

//func TestMediaCreate(t *testing.T) {
//	model.Init()
//	m := &model.Media{
//		Title:        "zzzz",
//		Introduction: "ewqewqe",
//	}
//	//s := &model.Sharpness{Model : model.Model{ID : 1}}
//	service.MediaCreate(m, []*model.MediaSharpness{
//		&model.MediaSharpness{SharpnessId: 1, Uri: "zzzz"},
//		&model.MediaSharpness{SharpnessId: 2, Uri: "zzzz"},
//	}, []*model.Category{
//		&model.Category{Model: model.Model{ID: 1}},
//	})
//
//	//model.DB().Where("user_id=?", user.ID).Delete(&model.UserRecord{})
//}

//func TestMediaUpdate(t *testing.T) {
//	model.Init()
//
//	//s := &model.Sharpness{Model : model.Model{ID : 1}}
//	if _, err := service.MediaUpdate(5, []*model.MediaSharpness{
//		&model.MediaSharpness{ID: 13, MediaID: 4, SharpnessId: 1, Uri: "qweqwe"},
//		&model.MediaSharpness{ID: 14, MediaID: 6, SharpnessId: 2, Uri: "eqweqw"},
//	}, []*model.Category{
//		&model.Category{Model: model.Model{ID: 2}},
//	}); err != nil {
//		t.Error(err)
//	}
//
//	//model.DB().Where("user_id=?", user.ID).Delete(&model.UserRecord{})
//}


func TestMSee(t *testing.T) {
	c := "a.txt"
	log.Println(filepath.Ext(c))
}

func TestUUID(t *testing.T) {
	for i := 0; i < 20; i++ {
		u := uuid.New()
		log.Println(u.String())
	}
}

func TestTOML(t *testing.T) {
	conf.Init()
	log.Println(conf.C())
}