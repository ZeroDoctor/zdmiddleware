package zdmiddleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

var limiterCache = cache.New(5*time.Minute, 10*time.Minute)

func NewRateLimiter(key func(*gin.Context) string, createLimiter func() (*rate.Limiter, time.Duration), abort func(*gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := key(c)

		limiter, ok := limiterCache.Get(k)
		if !ok {
			limiter, expire := createLimiter()
			limiterCache.Set(k, limiter, expire)
		}
		gin.Logger()

		if ok = limiter.(*rate.Limiter).Allow(); !ok {
			abort(c)
			return
		}

		c.Next()
	}
}
