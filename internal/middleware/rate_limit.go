package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware creates a rate limiting middleware
// Example: RateLimitMiddleware("10-M") allows 10 requests per minute
// Options: "5-S" (5 per second), "100-H" (100 per hour), "1000-D" (1000 per day)
func RateLimitMiddleware(rateFormat string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(rateFormat)
	if err != nil {
		panic("invalid rate limit format: " + err.Error())
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get client IP
		key := c.ClientIP()

		// Get limiter context
		context, err := instance.Get(c, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", string(rune(context.Limit)))
		c.Header("X-RateLimit-Remaining", string(rune(context.Remaining)))
		c.Header("X-RateLimit-Reset", string(rune(context.Reset)))

		// Check if limit exceeded
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": "too many requests, please try again later",
				"retry_after": context.Reset,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimitMiddleware applies stricter rate limiting for authentication endpoints
func AuthRateLimitMiddleware() gin.HandlerFunc {
	// 5 attempts per minute to prevent brute force attacks
	return RateLimitMiddleware("5-M")
}

// APIRateLimitMiddleware applies general rate limiting for API endpoints
func APIRateLimitMiddleware() gin.HandlerFunc {
	// 100 requests per minute for general API usage
	return RateLimitMiddleware("100-M")
}
