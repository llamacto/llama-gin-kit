package v1

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"github.com/llamacto/llama-gin-kit/pkg/database"
)

// RegisterRoutes registers all v1 version routes
func RegisterRoutes(engine *gin.Engine, v1 *gin.RouterGroup) {
	// Register health check routes
	RegisterHealthRoutes(v1)

	// Initialize repositories and services
	db := database.DB
	if db == nil {
		log.Fatal("Database connection not initialized")
	}

	// Initialize API key module
	apiKeyRepo := apikey.NewAPIKeyRepository(db)
	apiKeyService := apikey.NewAPIKeyService(apiKeyRepo)
	
	// Register API key routes
	RegisterAPIKeyRoutes(v1, apiKeyService)
}
