package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	v1 "github.com/zgiai/ginext/routes/v1"
)

// RegisterRoutes registers all routes
func RegisterRoutes(r *gin.Engine) {
	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API v1 routes
	v1Group := r.Group("/v1")
	v1.RegisterRoutes(r, v1Group)

	// API v2 routes will be added when needed
	// v2Group := r.Group("/v2")
}
