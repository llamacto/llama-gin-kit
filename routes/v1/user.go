package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/internal/modules/user"
	"github.com/llamacto/llama-gin-kit/pkg/database"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(v1 *gin.RouterGroup) {
	// Initialize user repository and service
	userRepo := user.NewUserRepository(database.GetDB())
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Public auth routes
	v1.POST("/register", userHandler.Register)
	v1.POST("/login", userHandler.Login)
	v1.POST("/password/reset", userHandler.ResetPassword)

	// User routes
	userGroup := v1.Group("/users")
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
}
