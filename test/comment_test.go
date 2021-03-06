package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/seachenjy/go-comment/config"
	"github.com/seachenjy/go-comment/dao"
	"golang.org/x/time/rate"
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

func TestTimeafter(t *testing.T) {
	t.Log(time.Now().UTC())
	t2 := time.Unix(time.Now().Unix()-5, 0).UTC()
	out := time.Now().Sub(t2)
	t.Log(out.Seconds())
}

func TestRate(t *testing.T) {
	l := rate.NewLimiter(20, 5)
	// c, _ := context.WithCancel(context.TODO())
	fmt.Println(l.Limit(), l.Burst())
	for i := 0; i < 100; i++ {
		if l.Allow() {
			fmt.Println(time.Now().Format("2016-01-02 15:04:05.000"))
		} else {
			fmt.Println("time out")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestAggregate(t *testing.T) {
	err := config.Init("../config.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", config.Cfg)

	mongo := dao.NewMongo(&config.Cfg)
	res := mongo.Aggregate("abc")

	t.Logf("%+v", res)
}
