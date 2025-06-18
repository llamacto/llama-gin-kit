package role

import (
	"time"
)

// Role represents a permission role within an organization
type Role struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
	Name           string     `gorm:"size:50;not null" json:"name"`
	DisplayName    string     `gorm:"size:100" json:"display_name"`
	Description    string     `gorm:"size:255" json:"description"`
	OrganizationID *uint      `json:"organization_id"` // If null, it's a system role
	// Permissions    string     `gorm:"type:json" json:"permissions"` // Temporarily disabled
	IsDefault bool `gorm:"default:false" json:"is_default"`
}

// TableName specifies the database table name
func (Role) TableName() string {
	return "organization_roles"
}

// Permission represents an individual permission that can be granted to a role
type Permission struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Name        string     `gorm:"size:100;not null;unique" json:"name"`
	DisplayName string     `gorm:"size:100" json:"display_name"`
	Description string     `gorm:"size:255" json:"description"`
	Category    string     `gorm:"size:50" json:"category"` // Grouping for UI display
}

// TableName specifies the database table name
func (Permission) TableName() string {
	return "organization_permissions"
}

// RoleWithPermissions combines role data with permissions for queries
type RoleWithPermissions struct {
	Role              Role         `json:"role"`
	PermissionDetails []Permission `json:"permission_details"`
}

// PermissionsByCategory groups permissions by category for easier UI display
type PermissionsByCategory struct {
	Category    string       `json:"category"`
	Permissions []Permission `json:"permissions"`
}
