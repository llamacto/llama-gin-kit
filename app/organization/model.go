package organization

import (
	"time"

	"github.com/llamacto/llama-gin-kit/app/user"
)

// Organization represents the organization model
type Organization struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	DisplayName string     `gorm:"size:100" json:"display_name"`
	Description string     `gorm:"size:500" json:"description"`
	Logo        string     `gorm:"size:255" json:"logo"`
	Website     string     `gorm:"size:255" json:"website"`
	Settings    string     `gorm:"type:json" json:"settings"` // JSON settings for organization
	Status      int        `gorm:"default:1" json:"status"`   // 1: active, 0: disabled
}

// TableName specifies the database table name
func (Organization) TableName() string {
	return "organizations"
}

// Team represents a team within an organization
type Team struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
	Name           string     `gorm:"size:100;not null" json:"name"`
	DisplayName    string     `gorm:"size:100" json:"display_name"`
	Description    string     `gorm:"size:500" json:"description"`
	OrganizationID uint       `gorm:"not null" json:"organization_id"`
	ParentTeamID   *uint      `json:"parent_team_id"` // For hierarchical team structure
	Settings       string     `gorm:"type:json" json:"settings"`
	Status         int        `gorm:"default:1" json:"status"` // 1: active, 0: disabled
}

// TableName specifies the database table name
func (Team) TableName() string {
	return "teams"
}

// Member represents a user's membership in an organization or team
type Member struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
	UserID         uint       `gorm:"not null" json:"user_id"`
	OrganizationID uint       `gorm:"not null" json:"organization_id"`
	TeamID         *uint      `json:"team_id"` // Optional, if member belongs to specific team
	RoleID         uint       `gorm:"not null" json:"role_id"`
	Status         int        `gorm:"default:1" json:"status"` // 1: active, 0: pending, 2: disabled
	JoinedAt       time.Time  `json:"joined_at"`
	InvitedBy      uint       `json:"invited_by"` // User ID who invited this member
}

// TableName specifies the database table name
func (Member) TableName() string {
	return "organization_members"
}

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
	Permissions    string     `gorm:"type:json" json:"permissions"`
	IsDefault      bool       `gorm:"default:false" json:"is_default"`
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

// Invitation represents a pending invitation to join an organization
type Invitation struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
	Email          string     `gorm:"size:100;not null" json:"email"`
	OrganizationID uint       `gorm:"not null" json:"organization_id"`
	TeamID         *uint      `json:"team_id"`
	RoleID         uint       `gorm:"not null" json:"role_id"`
	InvitedBy      uint       `json:"invited_by"`
	Token          string     `gorm:"size:100;not null" json:"token"`
	ExpiresAt      time.Time  `json:"expires_at"`
	Status         int        `gorm:"default:0" json:"status"` // 0: pending, 1: accepted, 2: rejected, 3: expired
}

// TableName specifies the database table name
func (Invitation) TableName() string {
	return "organization_invitations"
}

// OrganizationUser combines organization and user data for queries
type OrganizationUser struct {
	Organization Organization `json:"organization"`
	User         user.User    `json:"user"`
	Member       Member       `json:"member"`
	Role         Role         `json:"role"`
}

// TeamMember combines team and member data for queries
type TeamMember struct {
	Team   Team      `json:"team"`
	Member Member    `json:"member"`
	User   user.User `json:"user"`
	Role   Role      `json:"role"`
}
