package database

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/llamacto/llama-gin-kit/app/user"
	"gorm.io/gorm"
)

// RunMigrations runs all migrations for the application
func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migrations")
	startTime := time.Now()

	// Collect all migrations from different modules
	allMigrations := []*gormigrate.Migration{}

	// Add user migrations
	userMigrations := getUserMigrations()
	allMigrations = append(allMigrations, userMigrations...)

	// Add API key migrations (temporarily disabled)
	// apiKeyMigration := migrations.CreateAPIKeysTable()
	// allMigrations = append(allMigrations, apiKeyMigration)

	// Add organization migrations
	// orgMigrations := organization.GetMigrations()
	// allMigrations = append(allMigrations, orgMigrations...)

	// Initialize the migrator with all collected migrations
	m := gormigrate.New(db, &gormigrate.Options{
		TableName:      "migrations",
		IDColumnName:   "id",
		IDColumnSize:   255,
		UseTransaction: true,
	}, allMigrations)

	// Execute migrations
	if err := m.Migrate(); err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Printf("Migration completed successfully in %v", time.Since(startTime))
	return nil
}

// getUserMigrations returns migrations for the user module
func getUserMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202506180_create_users",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(&user.User{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("users")
			},
		},
		{
			ID: "202506181_create_default_users",
			Migrate: func(db *gorm.DB) error {
				// Create a default admin user if none exists
				var count int64
				db.Model(&user.User{}).Count(&count)

				if count == 0 {
					adminUser := &user.User{
						Username: "admin",
						Email:    "admin@example.com",
						Password: "hashed_password_here", // In a real app, this should be properly hashed
						Nickname: "Admin User",
						Status:   1, // 1: active, 0: disabled
					}

					result := db.Create(adminUser)
					return result.Error
				}

				return nil
			},
			Rollback: func(db *gorm.DB) error {
				return db.Where("username = ?", "admin").Delete(&user.User{}).Error
			},
		},
	}
}
