package database

import (
	"fmt"

	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"github.com/llamacto/llama-gin-kit/app/member"
	"github.com/llamacto/llama-gin-kit/app/organization"
	"github.com/llamacto/llama-gin-kit/app/team"
	"github.com/llamacto/llama-gin-kit/app/user"
	"github.com/llamacto/llama-gin-kit/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// getMigrations returns all migrations for the application
func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "20250620_initial_schema",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&user.User{},
					&organization.Organization{},
					&team.Team{},
					&apikey.APIKey{},
					&member.Member{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					&member.Member{},
					&apikey.APIKey{},
					&team.Team{},
					&organization.Organization{},
					&user.User{},
				)
			},
		},
	}
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

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
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

	// Run migrations
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	// Migrate the schema
	if err = m.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	DB = db
	return db, nil
}

// GetDB returns the database connection instance
func GetDB() *gorm.DB {
	return DB
}
