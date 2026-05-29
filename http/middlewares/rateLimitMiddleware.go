package middlewares

import (
	"sync"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	ipLimiters   = make(map[string]*rate.Limiter)
	userLimiters = make(map[string]*rate.Limiter)
	mu           sync.Mutex
)

func getIPLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := ipLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(2, 5)
		ipLimiters[ip] = limiter
	}
	return limiter
}

func getUserLimiter(userID string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := userLimiters[userID]
	if !exists {
		limiter = rate.NewLimiter(10, 20)
		userLimiters[userID] = limiter
	}
	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("sub")

		if exists {
			uid := userID.(string)
			limiter := getUserLimiter(uid)

			if !limiter.Allow() {
				helpers.Response(c, domain.HttpResponse{
					Code:    429,
					Message: "too many requests from your account, please try again later",
				})
				c.Abort()
				return
			}
		} else {
			ip := c.ClientIP()
			limiter := getIPLimiter(ip)

			if !limiter.Allow() {
				helpers.Response(c, domain.HttpResponse{
					Code:    429,
					Message: "too many requests from your IP address, please try again later",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
