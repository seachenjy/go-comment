package test

import (
	"fmt"
	"testing"

	"github.com/seachenjy/go-comment/config"
	"github.com/seachenjy/go-comment/dao"
)

func TestComment(t *testing.T) {
	c := dao.New()
	fmt.Println(c)
}

func TestConfig(t *testing.T) {
	err := config.Init("../config.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", config.Cfg)
}

func TestAddComment(t *testing.T) {
	err := config.Init("../config.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", config.Cfg)
	c := dao.New()
	c.SourceID = "abc"
	mongo := dao.NewMongo(&config.Cfg)
	if ok := c.Save(mongo); !ok {
		t.Error("save error")
	}
}

func TestGetComment(t *testing.T) {
	err := config.Init("../config.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", config.Cfg)

	mongo := dao.NewMongo(&config.Cfg)

	list := mongo.GetComments("abc", 0, 10)

	t.Logf("%++v", list[0])

}
