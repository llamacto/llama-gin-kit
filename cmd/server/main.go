package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/llamacto/llama-gin-kit/config"
	"github.com/llamacto/llama-gin-kit/pkg/database"
	"github.com/llamacto/llama-gin-kit/pkg/email"
	"github.com/llamacto/llama-gin-kit/pkg/jwt"
	"github.com/llamacto/llama-gin-kit/routes"
)

// @title Llama Gin Kit API
// @version 1.0
// @description A modern Go scaffold for AI-powered development with LLM integrations and agent-based architecture
// @host localhost:6066
// @BasePath /v1
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize JWT service
	jwt.Init(cfg)

	// Initialize email service
	email.Init(cfg)

	// Initialize database
	_, err = database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin mode
	gin.SetMode(gin.DebugMode)

	// Create Gin engine
	r := gin.Default()

	// Enable CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	// Register routes
	routes.RegisterRoutes(r)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "6066"
	}
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", serverAddr)

	go func() {
		if err := r.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
