package api

import (
	"fmt"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/seachenjy/go-comment/config"
	"github.com/seachenjy/go-comment/dao"
)

var (
	m      = &sync.Mutex{}
	inited bool
)

//Init init
func Init() {
	m.Lock()
	defer m.Unlock()
	if inited {
		return
	}
	inited = true
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	pprof.Register(r)
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.POST("/comment/save", saveComment)
	r.GET("/comments", comments)
	err := r.Run(fmt.Sprintf(`127.0.0.1:%d`, config.Cfg.Port))
	if err != nil {
		panic(err)
	}
}

//save comment api
func saveComment(c *gin.Context) {
	var parm dao.Comment
	c.BindJSON(&parm)
	c.JSON(200, parm)
}

//comments list api
func comments(c *gin.Context) {

}
