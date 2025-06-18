package apikey

import (
	"time"

	"gorm.io/gorm"
)

// APIKey represents an API key for authenticating API requests
type APIKey struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Key         string         `json:"key" gorm:"type:varchar(64);uniqueIndex;not null"` // Hashed key
	Prefix      string         `json:"prefix" gorm:"type:varchar(8);not null"`           // First 8 characters for identification
	UserID      uint           `json:"user_id" gorm:"not null"`                          // Owner of the API key
	LastUsedAt  *time.Time     `json:"last_used_at"`                                     // Track when the key was last used
	ExpiresAt   *time.Time     `json:"expires_at"`                                       // Optional expiration date
	Permissions string         `json:"permissions" gorm:"type:text"`                      // JSON string of permissions
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName specifies the table name for the APIKey model
func (APIKey) TableName() string {
	return "api_keys"
}
