package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/zgiai/ginext/pkg/logger"
	"github.com/zgiai/ginext/pkg/response"
)

// Recovery middleware handles panic recovery
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log stack trace
				logger.Error("Panic recovered", fmt.Errorf("%v", err))
				logger.Debug("Stack trace", string(debug.Stack()))

				response.Error(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
