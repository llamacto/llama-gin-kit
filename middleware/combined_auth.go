package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"github.com/llamacto/llama-gin-kit/pkg/middleware"
)

// CombinedAuth is a middleware that supports both API key and JWT authentication
// It will attempt to authenticate with API key first, then fall back to JWT if API key is not provided
func CombinedAuth(apiKeyService apikey.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for API key in header
		apiKeyHeader := c.GetHeader("X-API-Key")
		if apiKeyHeader == "" {
			// If no API key header, check in query parameter
			apiKeyHeader = c.Query("api_key")
		}
		
		// If API key is provided, use API key authentication
		if apiKeyHeader != "" {
			// Validate API key
			apiKeyObj, err := apiKeyService.ValidateAPIKey(apiKeyHeader)
			if err == nil {
				// API key is valid, set user ID and API key ID in context
				c.Set("userID", apiKeyObj.UserID)
				c.Set("apiKeyID", apiKeyObj.ID)
				c.Set("authType", "api_key")
				c.Next()
				return
			}
		}
		
		// If API key is not provided or is invalid, fall back to JWT auth
		jwtAuth := middleware.JWTAuth()
		jwtAuth(c)
		
		// If JWT auth was successful, set authType to jwt
		if !c.IsAborted() {
			c.Set("authType", "jwt")
		}
	}
}
