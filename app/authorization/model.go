package authorization

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a user role in the system
type Role struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string `gorm:"size:100;uniqueIndex;not null" json:"name"` // Role name (e.g., "admin", "user", "moderator")
	DisplayName string `gorm:"size:150;not null" json:"display_name"`     // Human readable name
	Description string `gorm:"type:text" json:"description"`              // Role description
	Level       int    `gorm:"default:0" json:"level"`                    // Role hierarchy level (higher = more permissions)
	IsSystem    bool   `gorm:"default:false" json:"is_system"`            // System roles cannot be deleted
	Status      int    `gorm:"default:1" json:"status"`                   // 1: active, 0: inactive

	// Relationships
	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []UserRole    `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

// Permission represents a specific permission in the system
type Permission struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string `gorm:"size:100;uniqueIndex;not null" json:"name"` // Permission name (e.g., "users.create")
	DisplayName string `gorm:"size:150;not null" json:"display_name"`     // Human readable name
	Description string `gorm:"type:text" json:"description"`              // Permission description
	Resource    string `gorm:"size:50;not null" json:"resource"`          // Resource type (e.g., "users", "organizations")
	Action      string `gorm:"size:50;not null" json:"action"`            // Action type (e.g., "create", "read", "update", "delete")
	Category    string `gorm:"size:50;default:'general'" json:"category"` // Permission category for grouping
	IsSystem    bool   `gorm:"default:false" json:"is_system"`            // System permissions cannot be deleted
	Status      int    `gorm:"default:1" json:"status"`                   // 1: active, 0: inactive

	// Relationships
	Roles []*Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// UserRole represents the relationship between users and roles
type UserRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID     uint       `gorm:"not null;index" json:"user_id"`
	RoleID     uint       `gorm:"not null;index" json:"role_id"`
	AssignedBy uint       `gorm:"index" json:"assigned_by"`      // User ID who assigned this role
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`          // Optional expiration date
	IsActive   bool       `gorm:"default:true" json:"is_active"` // Active status

	// Relationships
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// OrganizationRole represents organization-specific roles
type OrganizationRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID         uint `gorm:"not null;index" json:"user_id"`
	OrganizationID uint `gorm:"not null;index" json:"organization_id"`
	RoleID         uint `gorm:"not null;index" json:"role_id"`
	AssignedBy     uint `gorm:"index" json:"assigned_by"`
	IsActive       bool `gorm:"default:true" json:"is_active"`

	// Relationships
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TeamRole represents team-specific roles
type TeamRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID     uint `gorm:"not null;index" json:"user_id"`
	TeamID     uint `gorm:"not null;index" json:"team_id"`
	RoleID     uint `gorm:"not null;index" json:"role_id"`
	AssignedBy uint `gorm:"index" json:"assigned_by"`
	IsActive   bool `gorm:"default:true" json:"is_active"`

	// Relationships
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName methods for custom table names
func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (OrganizationRole) TableName() string {
	return "organization_roles"
}

func (TeamRole) TableName() string {
	return "team_roles"
}

// Policy represents a generic policy
type Policy struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Subject string `gorm:"size:100;not null" json:"subject"` // e.g., "role:1", "user:2"
	Action  string `gorm:"size:100;not null" json:"action"`  // e.g., "read", "write"
	Object  string `gorm:"size:100;not null" json:"object"`  // e.g., "article:1", "dataset:2"
	Effect  string `gorm:"size:10;not null" json:"effect"`   // "allow" or "deny"
}

// RolePermission is the explicit join table for the many-to-many relationship
// between Role and Permission.
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
	CreatedAt    time.Time
}

func (Policy) TableName() string {
	return "policies"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
