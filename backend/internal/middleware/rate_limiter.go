package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func RateLimiter(rps int, burst int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			limiter := rate.NewLimiter(rate.Limit(rps), burst)
			visitors[ip] = &visitor{limiter: limiter}
			v = visitors[ip]
		}
		mu.Unlock()

		if !v.limiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many request",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
