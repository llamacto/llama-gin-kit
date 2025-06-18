package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"github.com/llamacto/llama-gin-kit/pkg/middleware"
)

// RegisterAPIKeyRoutes registers routes related to API key management
func RegisterAPIKeyRoutes(v1 *gin.RouterGroup, apiKeyService apikey.Service) {
	// Create API key handler
	handler := apikey.NewAPIKeyHandler(apiKeyService)

	// API key management routes (needs JWT authentication)
	apikeyGroup := v1.Group("/apikeys")
	apikeyGroup.Use(middleware.JWTAuth())
	{
		apikeyGroup.POST("", handler.Create)
		apikeyGroup.GET("", handler.List)
		apikeyGroup.GET("/:id", handler.Get)
		apikeyGroup.PUT("/:id", handler.Update)
		apikeyGroup.DELETE("/:id", handler.Delete)
	}
}
