package database

import (
	"fmt"

	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/llamacto/llama-gin-kit/config"
	"github.com/llamacto/llama-gin-kit/pkg/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// getMigrations returns all migrations for the application
func getMigrations() []*gormigrate.Migration {
	allMigrations := []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				// This is a placeholder for the initial migration
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		// API Keys migration
		migrations.CreateAPIKeysTable(),
	}

	// Add organization migrations
	// allMigrations = append(allMigrations, organization.GetMigrations()...)

	return allMigrations
}

// InitDB initializes database connection and performs auto migration
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// Configure custom logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
		cfg.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(0) // Disable connection max lifetime

	// Check if we can connect to the database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Drop the table if it exists and recreate it with the correct structure
	err = db.Exec(`DROP TABLE IF EXISTS tts_audio_history`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to drop tts_audio_history table: %w", err)
	}

	// Create table manually with the structure matching our model
	err = db.Exec(`
		CREATE TABLE tts_audio_history (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			text TEXT NOT NULL,
			voice VARCHAR(50) DEFAULT 'alloy',
			audio_url VARCHAR(255) NOT NULL,
			duration NUMERIC(10,4),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create tts_audio_history table: %w", err)
	}

	// Create index on user_id
	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tts_audio_history_user_id ON tts_audio_history(user_id)`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create index on tts_audio_history: %w", err)
	}

	// Run migrations
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	// Migrate the schema
	if err = m.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// API keys table already migrated through gormigrate

	DB = db
	return db, nil
}

// GetDB returns the database connection instance
func GetDB() *gorm.DB {
	return DB
}
