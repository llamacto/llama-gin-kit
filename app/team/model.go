package team

import (
	"time"
)

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
	// Settings       string     `gorm:"type:json;default:'{}'" json:"settings"` // Temporarily disabled
	Status int `gorm:"default:1" json:"status"` // 1: active, 0: disabled
}

// TableName specifies the database table name
func (Team) TableName() string {
	return "teams"
}

// TeamWithStats includes team data with member statistics
type TeamWithStats struct {
	Team        Team  `json:"team"`
	MemberCount int64 `json:"member_count"`
}

// TeamHierarchy represents a team with its parent and children information
type TeamHierarchy struct {
	Team     Team   `json:"team"`
	Parent   *Team  `json:"parent,omitempty"`
	Children []Team `json:"children,omitempty"`
}
