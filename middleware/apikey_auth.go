package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/apikey"
)

// APIKeyAuth is a middleware for API key authentication
func APIKeyAuth(apiKeyService apikey.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for API key in header
		apiKeyHeader := c.GetHeader("X-API-Key")
		
		// If no API key in header, check for it in query parameters
		if apiKeyHeader == "" {
			apiKeyHeader = c.Query("api_key")
		}
		
		// If still no API key, return error
		if apiKeyHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "API key is required",
			})
			c.Abort()
			return
		}
		
		// Validate API key
		apiKeyObj, err := apiKeyService.ValidateAPIKey(apiKeyHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid API key",
			})
			c.Abort()
			return
		}
		
		// Store user ID and API key ID in context
		c.Set("userID", apiKeyObj.UserID)
		c.Set("apiKeyID", apiKeyObj.ID)
		
		// If specific permissions are required, check them
		if requiredPerms, exists := c.Get("requiredPermissions"); exists {
			if !hasPermissions(apiKeyObj.Permissions, requiredPerms.([]string)) {
				c.JSON(http.StatusForbidden, gin.H{
					"code": 403,
					"msg":  "API key does not have required permissions",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// RequirePermissions is a middleware for requiring specific permissions
func RequirePermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("requiredPermissions", permissions)
		c.Next()
	}
}

// hasPermissions checks if the API key has the required permissions
func hasPermissions(apiKeyPerms string, requiredPerms []string) bool {
	// If no permissions are required, allow access
	if len(requiredPerms) == 0 {
		return true
	}
	
	// If API key has no permissions, deny access
	if apiKeyPerms == "" {
		return false
	}
	
	// Split API key permissions
	perms := strings.Split(apiKeyPerms, ",")
	
	// Check if API key has all required permissions
	permMap := make(map[string]bool)
	for _, p := range perms {
		permMap[strings.TrimSpace(p)] = true
	}
	
	// Check for wildcard permission
	if permMap["*"] {
		return true
	}
	
	// Check for each required permission
	for _, required := range requiredPerms {
		if !permMap[required] {
			return false
		}
	}
	
	return true
}
