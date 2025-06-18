package invitation

import (
	"time"
)

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

// InvitationWithDetails combines invitation data with related entities for queries
type InvitationWithDetails struct {
	ID               uint      `json:"id"`
	Email            string    `json:"email"`
	OrganizationID   uint      `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	TeamID           *uint     `json:"team_id"`
	TeamName         *string   `json:"team_name"`
	RoleID           uint      `json:"role_id"`
	RoleName         string    `json:"role_name"`
	RoleDisplayName  string    `json:"role_display_name"`
	InvitedBy        uint      `json:"invited_by"`
	InviterName      string    `json:"inviter_name"`
	InviterEmail     string    `json:"inviter_email"`
	Token            string    `json:"token"`
	ExpiresAt        time.Time `json:"expires_at"`
	Status           int       `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// InvitationStats represents invitation statistics
type InvitationStats struct {
	Total    int64 `json:"total"`
	Pending  int64 `json:"pending"`
	Accepted int64 `json:"accepted"`
	Rejected int64 `json:"rejected"`
	Expired  int64 `json:"expired"`
}
