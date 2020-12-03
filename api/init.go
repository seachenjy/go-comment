package api

import (
	"fmt"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/seachenjy/go-comment/config"
)

var (
	m      = &sync.Mutex{}
	inited bool
)

//Init init
func Init(cfg *config.Config) {
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
	r.GET("/comment/save", func(c *gin.Context) {

	})
	r.GET("/comments", func(c *gin.Context) {

	})
	err := r.Run(fmt.Sprintf(`127.0.0.1:%d`, cfg.Port))
	if err != nil {
		panic(err)
	}
}
