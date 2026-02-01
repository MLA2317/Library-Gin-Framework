package middleware

import (
	"library-project/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Log the request
		fields := map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": statusCode,
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}

		// Add user ID if authenticated
		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = userID
		}

		// Log based on status code
		if statusCode >= 500 {
			utils.LogError(nil, "HTTP request completed with server error", fields)
		} else if statusCode >= 400 {
			utils.LogWarning("HTTP request completed with client error", fields)
		} else {
			utils.LogInfo("HTTP request completed", fields)
		}
	}
}
