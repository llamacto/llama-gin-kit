package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"github.com/llamacto/llama-gin-kit/app/organization"
	apikeyMiddleware "github.com/llamacto/llama-gin-kit/middleware"
)

// RegisterOrganizationRoutes registers organization routes
func RegisterOrganizationRoutes(router *gin.RouterGroup, handler *organization.Handler, apiKeyService apikey.Service) {
	// Routes that require authentication
	authRouter := router.Group("")
	authRouter.Use(apikeyMiddleware.CombinedAuth(apiKeyService))

	// Organization endpoints - only core organization functionality
	orgRouter := authRouter.Group("/organizations")
	orgRouter.POST("", handler.CreateOrganization)
	orgRouter.GET("", handler.ListOrganizations)
	orgRouter.GET("/me", handler.GetMyOrganizations)
	orgRouter.GET("/:id", handler.GetOrganization)
	orgRouter.PUT("/:id", handler.UpdateOrganization)
	orgRouter.DELETE("/:id", handler.DeleteOrganization)
}
