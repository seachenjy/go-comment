package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/seachenjy/go-comment/config"
	"github.com/seachenjy/go-comment/dao"
	"github.com/seachenjy/go-comment/log"
	"github.com/sirupsen/logrus"
)

var (
	m      = &sync.Mutex{}
	inited bool
	d      dao.Dao

	//APIErrors all api error code and message
	APIErrors = map[int]string{
		1001: "params decode error",
		1002: "content too short",
		1003: "source id can't be empty",
		1004: "save comment error",
		1005: "timeout",
	}
)

//Init init
func Init() {
	m.Lock()
	defer m.Unlock()
	if inited {
		return
	}
	inited = true
	if config.Cfg.Db == "mongo" {
		d = dao.NewMongo(&config.Cfg)
	}
	go ipLimitTask()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	pprof.Register(r)
	r.Use(logger)
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.POST("/comment/save", saveComment)
	r.GET("/comments", comments)

	err := r.Run(fmt.Sprintf(`0.0.0.0:%d`, config.Cfg.Port))
	if err != nil {
		log.GetLogger().Error(err)
	}
}

//save comment api
func saveComment(c *gin.Context) {
	if !iplimit(c) {
		return
	}
	comment := dao.New()
	if err := c.BindJSON(comment); err != nil {
		throwError(1001, c)
		return
	}
	if comment.Content == "" || len(comment.Content) <= 5 {
		throwError(1002, c)
		return
	}
	if comment.SourceID == "" {
		throwError(1003, c)
		return
	}
	comment.IPAddress = c.ClientIP()

	if ok := comment.Save(d); !ok {
		throwError(1004, c)
		return
	}
	c.JSON(200, comment)
}

type commentsAPIParams struct {
	dao.SourceID `form:"source_id"`
	Offset       int64 `form:"offset"`
	Limit        int64 `form:"limit"`
}

//comments list api
func comments(c *gin.Context) {
	apiparams := &commentsAPIParams{}
	if err := c.BindQuery(apiparams); err != nil {
		throwError(1001, c)
		return
	}
	if apiparams.SourceID == "" {
		throwError(1003, c)
		return
	}
	comments, aggregate := dao.Get(apiparams.SourceID, d, apiparams.Offset, apiparams.Limit)

	c.JSON(200, gin.H{
		"data": gin.H{
			"comments":  comments,
			"aggregate": aggregate,
		},
	})
}

func throwError(code int, c *gin.Context) {
	log.GetLogger().Error(APIErrors[code])
	c.JSON(500, gin.H{
		"status": 0,
		"code":   code,
	})
}

func logger(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)
	reqMethod := c.Request.Method
	reqURL := c.Request.RequestURI
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()

	log.GetLogger().WithFields(logrus.Fields{
		"status_code":  statusCode,
		"latency_time": latencyTime,
		"client_ip":    clientIP,
		"req_method":   reqMethod,
		"req_uri":      reqURL,
		"header":       c.Request.Header,
	}).Info()
}
