package v1

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zgiai/ginext/internal/modules/explain"
	"github.com/zgiai/ginext/internal/modules/tts"
	"github.com/zgiai/ginext/internal/modules/user"
	"github.com/zgiai/ginext/pkg/database"
	"github.com/zgiai/ginext/pkg/middleware"
)

// RegisterRoutes registers all v1 version routes
func RegisterRoutes(engine *gin.Engine, v1 *gin.RouterGroup) {
	// Register health check routes
	RegisterHealthRoutes(v1)

	// Initialize repositories and services
	db := database.GetDB()
	if db == nil {
		log.Fatal("Database connection not initialized")
	}

	// Initialize user module
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Public auth routes
	v1.POST("/register", userHandler.Register)
	v1.POST("/login", userHandler.Login)
	v1.POST("/password/reset", userHandler.ResetPassword)

	// Protected user routes
	userGroup := v1.Group("/users")
	userGroup.Use(middleware.JWTAuth())
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

	// Initialize TTS service for explain module
	ttsRepo := tts.NewTTSRepository(db)
	ttsService := tts.NewTTSService(ttsRepo)

	// Initialize Explain module
	if err := explain.InitModule(engine, db, ttsService); err != nil {
		log.Printf("Failed to initialize explain module: %v", err)
	}

	// Register TTS routes
	RegisterTTSRoutes(v1)
}
