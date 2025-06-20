package member

import (
	"time"

	"github.com/llamacto/llama-gin-kit/app/organization"
	"github.com/llamacto/llama-gin-kit/app/user"
	"gorm.io/gorm"
)

// Member represents a user's membership in an organization or team
type Member struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserID         uint           `gorm:"not null" json:"user_id"`
	OrganizationID uint           `gorm:"not null" json:"organization_id"`
	TeamID         *uint          `json:"team_id"`                 // Pointer to allow null
	Status         int            `gorm:"default:1" json:"status"` // 1: active, 0: pending, 2: disabled
	JoinedAt       time.Time      `json:"joined_at"`
	InvitedBy      uint           `json:"invited_by"` // User ID who invited this member

	// Relationships
	User         user.User                 `gorm:"foreignKey:UserID"`
	Organization organization.Organization `gorm:"foreignKey:OrganizationID"`
}

// TableName specifies the database table name
func (Member) TableName() string {
	return "organization_members"
}

// MemberWithDetails combines member data with related entities for queries
type MemberWithDetails struct {
	ID               uint      `json:"id"`
	UserID           uint      `json:"user_id"`
	UserName         string    `json:"user_name"`
	UserEmail        string    `json:"user_email"`
	UserNickname     string    `json:"user_nickname"`
	UserAvatar       string    `json:"user_avatar"`
	OrganizationID   uint      `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	TeamID           *uint     `json:"team_id"`
	TeamName         *string   `json:"team_name"`
	RoleID           uint      `json:"role_id"`
	RoleName         string    `json:"role_name"`
	RoleDisplayName  string    `json:"role_display_name"`
	Status           int       `json:"status"`
	JoinedAt         time.Time `json:"joined_at"`
	InvitedBy        uint      `json:"invited_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// OrganizationMember combines organization and member data for queries
type OrganizationMember struct {
	Member Member    `json:"member"`
	User   user.User `json:"user"`
}

// TeamMember combines team and member data for queries
type TeamMember struct {
	Member Member    `json:"member"`
	User   user.User `json:"user"`
}
