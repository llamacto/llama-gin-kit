package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/llamacto/llama-gin-kit/routes/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ProjectInfo represents the project information response
type ProjectInfo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	GoVersion   string  `json:"go_version"`
	BuildTime   string  `json:"build_time"`
	Environment string  `json:"environment"`
	API         APIInfo `json:"api"`
	Links       Links   `json:"links"`
}

// APIInfo represents API information
type APIInfo struct {
	Version   string   `json:"version"`
	BaseURL   string   `json:"base_url"`
	Endpoints []string `json:"endpoints"`
	Features  []string `json:"features"`
}

// Links represents useful links
type Links struct {
	Documentation string `json:"documentation"`
	Health        string `json:"health"`
	Swagger       string `json:"swagger"`
}

// RegisterRoutes registers all routes
func RegisterRoutes(r *gin.Engine) {
	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint - Project information
	r.GET("/", func(c *gin.Context) {
		info := ProjectInfo{
			Name:        "Llama Gin Kit",
			Description: "A comprehensive REST API service built with Go and Gin framework, featuring user management, organization/team management, API key authentication, and modular architecture.",
			Version:     "v1.0.0",
			GoVersion:   "1.23.0+",
			BuildTime:   time.Now().Format("2006-01-02 15:04:05"),
			Environment: gin.Mode(),
			API: APIInfo{
				Version: "v1",
				BaseURL: "/v1",
				Endpoints: []string{
					"POST /v1/register - User registration",
					"POST /v1/login - User login",
					"GET /v1/users/profile - Get user profile",
					"POST /v1/organizations - Create organization",
					"GET /v1/organizations - List organizations",
					"POST /v1/teams - Create team",
					"GET /v1/teams/:id - Get team details",
					"POST /v1/apikeys - Create API key",
					"GET /v1/apikeys - List API keys",
				},
				Features: []string{
					"JWT Authentication",
					"API Key Authentication",
					"User Management",
					"Organization Management",
					"Team Management",
					"Role-based Access Control",
					"Email Notifications",
					"PostgreSQL Database",
					"Docker Support",
					"Swagger Documentation",
				},
			},
			Links: Links{
				Documentation: "/swagger/index.html",
				Health:        "/v1/health/status",
				Swagger:       "/swagger/*any",
			},
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Welcome to Llama Gin Kit API",
			"data":    info,
		})
	})

	// Legacy ping endpoint
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
