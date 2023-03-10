package zdmiddleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

type GinRateLimiter struct {
	KeyFunc       func(*gin.Context) string
	CreateLimiter func() (*rate.Limiter, time.Duration)
	Abort         func(*gin.Context)

	defaultExpiration time.Duration
	cleanupExpiration time.Duration

	limiterCache *cache.Cache
}

func NewReateLimiter(defaultExpiration time.Duration, cleanupExpiration time.Duration) *GinRateLimiter {
	return &GinRateLimiter{
		defaultExpiration: defaultExpiration,
		cleanupExpiration: cleanupExpiration,
		limiterCache:      cache.New(defaultExpiration, cleanupExpiration),
	}
}

func (g *GinRateLimiter) Limiter(c *gin.Context) {
	key := g.KeyFunc(c)

	limiter, ok := g.limiterCache.Get(key)
	if !ok {
		limiter, expire := g.CreateLimiter()
		g.limiterCache.Set(key, limiter, expire)
	}

	if ok = limiter.(*rate.Limiter).Allow(); !ok {
		g.Abort(c)
		return
	}

	c.Next()
}
