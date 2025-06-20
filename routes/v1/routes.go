package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"github.com/llamacto/llama-gin-kit/app/authorization"
	"github.com/llamacto/llama-gin-kit/app/organization"
	"github.com/llamacto/llama-gin-kit/app/user"
	"github.com/llamacto/llama-gin-kit/middleware"
	"github.com/llamacto/llama-gin-kit/pkg/database"
	pkgmiddleware "github.com/llamacto/llama-gin-kit/pkg/middleware"
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

	// Initialize user module
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Register user routes
	// Public auth routes
	v1.POST("/register", userHandler.Register)
	v1.POST("/login", userHandler.Login)
	v1.POST("/password/reset", userHandler.ResetPassword)

	// Protected user routes
	userGroup := v1.Group("/users")
	userGroup.Use(pkgmiddleware.JWTAuth())
	{
		userGroup.GET("/profile", userHandler.GetProfile)
		userGroup.PUT("/profile", userHandler.UpdateProfile)
		userGroup.PUT("/password", userHandler.ChangePassword)
		userGroup.DELETE("/account", userHandler.DeleteAccount)

		// Admin routes
		userGroup.GET("", userHandler.List)
		userGroup.GET("/:id", userHandler.Get)
		userGroup.GET("/:id/info", userHandler.GetUserInfo)
	}

	// Initialize API key module
	apiKeyRepo := apikey.NewAPIKeyRepository(db)
	apiKeyService := apikey.NewAPIKeyService(apiKeyRepo)

	// Register API key routes
	RegisterAPIKeyRoutes(v1, apiKeyService)

	// Initialize organization module
	orgRepo := organization.NewRepository(db)
	orgService := organization.NewService(orgRepo, userService, db)
	orgHandler := organization.NewHandler(orgService)

	// Register organization routes
	RegisterOrganizationRoutes(v1, orgHandler, apiKeyService)

	// Register team routes
	TeamRoutes(v1)

	// Initialize authorization components
	authRepo := authorization.NewRepository(db)
	authSvc := authorization.NewService(authRepo)

	// Register auth routes
	RegisterAuthRoutes(v1, authSvc)

	// Example of a route that accepts either JWT or API key authentication
	// 使用CombinedAuth中间件，支持JWT和API key双重认证
	combinedAuthMiddleware := middleware.CombinedAuth(apiKeyService)
	v1.GET("/protected", combinedAuthMiddleware, func(c *gin.Context) {
		// 获取认证类型
		authType := c.GetString("authType")
		userID := c.GetUint("userID")

		c.JSON(http.StatusOK, gin.H{
			"message":   "认证成功",
			"auth_type": authType,
			"user_id":   userID,
		})
	})
}
