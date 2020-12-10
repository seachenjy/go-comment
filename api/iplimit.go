package api

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var allowIPPool = &sync.Map{}

//ipLimitTask clear cache
func ipLimitTask() {
	l := rate.NewLimiter(20, 5)
	for {
		if l.Allow() {
			allowIPPool.Range(func(key, value interface{}) bool {
				ip := key.(string)
				age := value.(time.Time)
				if time.Now().Sub(age).Seconds() > 5 {
					allowIPPool.Delete(ip)
				}
				return true
			})
		}
	}
}

func iplimit(c *gin.Context) bool {
	ip := c.ClientIP()
	if age, ok := allowIPPool.Load(ip); ok {
		age := age.(time.Time)
		if time.Now().Sub(age).Seconds() < 5 {
			allowIPPool.Store(ip, time.Now())
			throwError(1005, c)
			return false
		}
		allowIPPool.Store(ip, time.Now())
		return true

	}
	allowIPPool.Store(ip, time.Now())
	return true

}
