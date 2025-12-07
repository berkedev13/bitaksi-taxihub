package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ClientLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	burst    int
}

func NewClientLimiter(r rate.Limit, burst int) *ClientLimiter {
	return &ClientLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		burst:    burst,
	}
}

func (cl *ClientLimiter) getLimiter(ip string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	limiter, exists := cl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(cl.r, cl.burst)
		cl.limiters[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware(cl *ClientLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := cl.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}
