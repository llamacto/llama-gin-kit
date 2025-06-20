package organization

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// JSONString is a custom type for handling JSON strings in GORM
type JSONString string

// Value implements the driver.Valuer interface
func (j JSONString) Value() (driver.Value, error) {
	if j == "" {
		return "{}", nil
	}
	return string(j), nil
}

// Scan implements the sql.Scanner interface
func (j *JSONString) Scan(value interface{}) error {
	if value == nil {
		*j = "{}"
		return nil
	}
	switch v := value.(type) {
	case string:
		*j = JSONString(v)
	case []byte:
		*j = JSONString(v)
	default:
		return fmt.Errorf("cannot scan %T into JSONString", value)
	}
	return nil
}

// Organization represents the organization model
type Organization struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"size:500" json:"description"`
	Logo        string         `gorm:"size:255" json:"logo"`
	Website     string         `gorm:"size:255" json:"website"`
	// Settings    *string        `gorm:"type:json" json:"settings,omitempty"` // JSON settings for organization - temporarily disabled
	Status int `gorm:"default:1" json:"status"` // 1: active, 0: disabled
}

// TableName specifies the database table name
func (Organization) TableName() string {
	return "organizations"
}

// OrganizationStats includes organization data with statistics
type OrganizationStats struct {
	Organization Organization `json:"organization"`
	MemberCount  int64        `json:"member_count"`
	TeamCount    int64        `json:"team_count"`
	RoleCount    int64        `json:"role_count"`
}
