package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/pkg/logger"
)

// Logger creates a middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Log request details
		logger.Info(fmt.Sprintf(
			"Request: method=%s path=%s status=%d latency=%v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		))
	}
}
