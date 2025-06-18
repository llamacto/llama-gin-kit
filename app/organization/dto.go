package organization

import (
	"time"
)

// CreateOrganizationRequest represents the request to create an organization
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Website     string `json:"website"`
	Settings    string `json:"settings,omitempty"`
}

// UpdateOrganizationRequest represents the request to update an organization
type UpdateOrganizationRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Website     string `json:"website"`
	Settings    string `json:"settings,omitempty"`
	Status      *int   `json:"status,omitempty"`
}

// OrganizationResponse represents the organization data in responses
type OrganizationResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	Website     string    `json:"website"`
	Settings    string    `json:"settings,omitempty"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrganizationStatsResponse represents organization statistics
type OrganizationStatsResponse struct {
	Organization OrganizationResponse `json:"organization"`
	MemberCount  int64                `json:"member_count"`
	TeamCount    int64                `json:"team_count"`
	RoleCount    int64                `json:"role_count"`
}

// CreateTeamRequest represents the request to create a team
type CreateTeamRequest struct {
	Name           string `json:"name" binding:"required"`
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	ParentTeamID   *uint  `json:"parent_team_id,omitempty"`
	Settings       string `json:"settings,omitempty"`
}

// UpdateTeamRequest represents the request to update a team
type UpdateTeamRequest struct {
	DisplayName  string `json:"display_name"`
	Description  string `json:"description"`
	ParentTeamID *uint  `json:"parent_team_id,omitempty"`
	Settings     string `json:"settings,omitempty"`
	Status       *int   `json:"status,omitempty"`
}

// TeamResponse represents the team data in responses
type TeamResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	DisplayName    string    `json:"display_name"`
	Description    string    `json:"description"`
	OrganizationID uint      `json:"organization_id"`
	ParentTeamID   *uint     `json:"parent_team_id,omitempty"`
	Settings       string    `json:"settings,omitempty"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AddMemberRequest represents the request to add a member
type AddMemberRequest struct {
	UserID         uint  `json:"user_id" binding:"required"`
	OrganizationID uint  `json:"organization_id" binding:"required"`
	TeamID         *uint `json:"team_id,omitempty"`
	RoleID         uint  `json:"role_id" binding:"required"`
}

// UpdateMemberRequest represents the request to update a member
type UpdateMemberRequest struct {
	TeamID *uint `json:"team_id,omitempty"`
	RoleID uint  `json:"role_id" binding:"required"`
	Status *int  `json:"status,omitempty"`
}

// MemberResponse represents the member data in responses
type MemberResponse struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"user_id"`
	OrganizationID uint      `json:"organization_id"`
	TeamID         *uint     `json:"team_id,omitempty"`
	RoleID         uint      `json:"role_id"`
	Status         int       `json:"status"`
	JoinedAt       time.Time `json:"joined_at"`
	InvitedBy      uint      `json:"invited_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateRoleRequest represents the request to create a role
type CreateRoleRequest struct {
	Name           string `json:"name" binding:"required"`
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
	OrganizationID *uint  `json:"organization_id,omitempty"`
	Permissions    string `json:"permissions" binding:"required"`
	IsDefault      bool   `json:"is_default"`
}

// UpdateRoleRequest represents the request to update a role
type UpdateRoleRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Permissions string `json:"permissions"`
	IsDefault   *bool  `json:"is_default,omitempty"`
}

// RoleResponse represents the role data in responses
type RoleResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	DisplayName    string    `json:"display_name"`
	Description    string    `json:"description"`
	OrganizationID *uint     `json:"organization_id,omitempty"`
	Permissions    string    `json:"permissions"`
	IsDefault      bool      `json:"is_default"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateInvitationRequest represents the request to create an invitation
type CreateInvitationRequest struct {
	Email          string `json:"email" binding:"required,email"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	TeamID         *uint  `json:"team_id,omitempty"`
	RoleID         uint   `json:"role_id" binding:"required"`
}

// InvitationResponse represents the invitation data in responses
type InvitationResponse struct {
	ID             uint      `json:"id"`
	Email          string    `json:"email"`
	OrganizationID uint      `json:"organization_id"`
	TeamID         *uint     `json:"team_id,omitempty"`
	RoleID         uint      `json:"role_id"`
	InvitedBy      uint      `json:"invited_by"`
	Token          string    `json:"token,omitempty"`
	ExpiresAt      time.Time `json:"expires_at"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AcceptInvitationRequest represents the request to accept an invitation
type AcceptInvitationRequest struct {
	Token string `json:"token" binding:"required"`
}

// CheckPermissionRequest represents the request to check a permission
type CheckPermissionRequest struct {
	OrganizationID uint   `json:"organization_id" binding:"required"`
	Permission     string `json:"permission" binding:"required"`
}

// CheckPermissionResponse represents the response from a permission check
type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Data  interface{} `json:"data"`
}
