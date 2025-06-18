package main

import (
	"log"

	"github.com/llamacto/llama-gin-kit/config"
	"github.com/llamacto/llama-gin-kit/pkg/database"
)

func main() {
	log.Println("Starting database migration process")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations using the new migration system
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migration completed successfully")
}
