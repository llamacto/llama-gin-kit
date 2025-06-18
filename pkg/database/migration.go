package database

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/llamacto/llama-gin-kit/app/organization"
	"github.com/llamacto/llama-gin-kit/app/user"
	"gorm.io/gorm"
)

// RunMigrations runs all migrations for the application
func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migrations")
	startTime := time.Now()

	// Collect all migrations from different modules
	migrations := []*gormigrate.Migration{}
	
	// Add user migrations
	userMigrations := getUserMigrations()
	migrations = append(migrations, userMigrations...)
	
	// Add organization migrations
	orgMigrations := organization.GetMigrations()
	migrations = append(migrations, orgMigrations...)
	
	// Initialize the migrator with all collected migrations
	m := gormigrate.New(db, &gormigrate.Options{
		TableName:      "migrations",
		IDColumnName:   "id",
		IDColumnSize:   255,
		UseTransaction: true,
	}, migrations)

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
