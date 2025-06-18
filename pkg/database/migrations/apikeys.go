package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"gorm.io/gorm"
)

// CreateAPIKeysTable creates the api_keys table
func CreateAPIKeysTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202506181130_create_api_keys_table",
		Migrate: func(tx *gorm.DB) error {
			// Create the table
			return tx.AutoMigrate(&apikey.APIKey{})
		},
		Rollback: func(tx *gorm.DB) error {
			// Drop the table
			return tx.Migrator().DropTable("api_keys")
		},
	}
}
