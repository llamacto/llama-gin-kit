package authorization

import "time"

// ===== Role DTOs =====

// CreateRoleRequest represents the request to create a role
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100" example:"admin"`
	DisplayName string `json:"display_name" binding:"required,min=3,max=150" example:"Administrator"`
	Description string `json:"description" example:"System administrator with full access"`
	Level       int    `json:"level" example:"100"`
	Status      int    `json:"status" example:"1"`
}

// UpdateRoleRequest represents the request to update a role
type UpdateRoleRequest struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description,omitempty"`
	Level       *int    `json:"level,omitempty"`
	Status      *int    `json:"status,omitempty"`
}

// RoleResponse represents the role response
type RoleResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Level       int       `json:"level"`
	IsSystem    bool      `json:"is_system"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoleWithPermissionsResponse represents role with its permissions
type RoleWithPermissionsResponse struct {
	RoleResponse
	Permissions []PermissionResponse `json:"permissions"`
}

// ===== Permission DTOs =====

// CreatePermissionRequest represents the request to create a permission
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100" example:"users.create"`
	DisplayName string `json:"display_name" binding:"required,min=3,max=150" example:"Create Users"`
	Description string `json:"description" example:"Permission to create new users"`
	Resource    string `json:"resource" binding:"required,min=2,max=50" example:"users"`
	Action      string `json:"action" binding:"required,min=2,max=50" example:"create"`
	Category    string `json:"category" example:"user_management"`
	Status      int    `json:"status" example:"1"`
}

// UpdatePermissionRequest represents the request to update a permission
type UpdatePermissionRequest struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description,omitempty"`
	Resource    *string `json:"resource,omitempty" binding:"omitempty,min=2,max=50"`
	Action      *string `json:"action,omitempty" binding:"omitempty,min=2,max=50"`
	Category    *string `json:"category,omitempty"`
	Status      *int    `json:"status,omitempty"`
}

// PermissionResponse represents the permission response
type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Category    string    `json:"category"`
	IsSystem    bool      `json:"is_system"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ===== Role Permission Assignment DTOs =====

// AssignPermissionsRequest represents the request to assign permissions to a role
type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required,min=1" example:"[1,2,3]"`
}

// RemovePermissionsRequest represents the request to remove permissions from a role
type RemovePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required,min=1" example:"[1,2,3]"`
}

// ===== User Role Assignment DTOs =====

// AssignRoleRequest represents the request to assign a role to a user
type AssignRoleRequest struct {
	UserID    uint       `json:"user_id" binding:"required" example:"1"`
	RoleID    uint       `json:"role_id" binding:"required" example:"1"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// AssignRolesRequest represents the request to assign multiple roles to a user
type AssignRolesRequest struct {
	UserID  uint   `json:"user_id" binding:"required" example:"1"`
	RoleIDs []uint `json:"role_ids" binding:"required,min=1" example:"[1,2,3]"`
}

// UserRoleResponse represents the user role assignment response
type UserRoleResponse struct {
	ID         uint         `json:"id"`
	UserID     uint         `json:"user_id"`
	RoleID     uint         `json:"role_id"`
	AssignedBy uint         `json:"assigned_by"`
	ExpiresAt  *time.Time   `json:"expires_at,omitempty"`
	IsActive   bool         `json:"is_active"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	Role       RoleResponse `json:"role"`
}

// ===== Organization Role DTOs =====

// AssignOrganizationRoleRequest represents the request to assign a role within an organization
type AssignOrganizationRoleRequest struct {
	UserID         uint `json:"user_id" binding:"required" example:"1"`
	OrganizationID uint `json:"organization_id" binding:"required" example:"1"`
	RoleID         uint `json:"role_id" binding:"required" example:"1"`
}

// OrganizationRoleResponse represents the organization role assignment response
type OrganizationRoleResponse struct {
	ID             uint         `json:"id"`
	UserID         uint         `json:"user_id"`
	OrganizationID uint         `json:"organization_id"`
	RoleID         uint         `json:"role_id"`
	AssignedBy     uint         `json:"assigned_by"`
	IsActive       bool         `json:"is_active"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Role           RoleResponse `json:"role"`
}

// ===== Team Role DTOs =====

// AssignTeamRoleRequest represents the request to assign a role within a team
type AssignTeamRoleRequest struct {
	UserID uint `json:"user_id" binding:"required" example:"1"`
	TeamID uint `json:"team_id" binding:"required" example:"1"`
	RoleID uint `json:"role_id" binding:"required" example:"1"`
}

// TeamRoleResponse represents the team role assignment response
type TeamRoleResponse struct {
	ID         uint         `json:"id"`
	UserID     uint         `json:"user_id"`
	TeamID     uint         `json:"team_id"`
	RoleID     uint         `json:"role_id"`
	AssignedBy uint         `json:"assigned_by"`
	IsActive   bool         `json:"is_active"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	Role       RoleResponse `json:"role"`
}

// ===== Policy DTOs =====

// CreatePolicyRequest represents the request to create a policy
type CreatePolicyRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100" example:"admin_access"`
	DisplayName string `json:"display_name" binding:"required,min=3,max=150" example:"Admin Access Policy"`
	Description string `json:"description" example:"Policy for admin access control"`
	Resource    string `json:"resource" binding:"required,min=2,max=50" example:"users"`
	Effect      string `json:"effect" binding:"required,oneof=allow deny" example:"allow"`
	Conditions  string `json:"conditions" example:"{\"time_range\": \"9:00-18:00\"}"`
	Priority    int    `json:"priority" example:"100"`
}

// UpdatePolicyRequest represents the request to update a policy
type UpdatePolicyRequest struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description,omitempty"`
	Resource    *string `json:"resource,omitempty" binding:"omitempty,min=2,max=50"`
	Effect      *string `json:"effect,omitempty" binding:"omitempty,oneof=allow deny"`
	Conditions  *string `json:"conditions,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// PolicyResponse represents the policy response
type PolicyResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Effect      string    `json:"effect"`
	Conditions  string    `json:"conditions"`
	Priority    int       `json:"priority"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ===== Permission Check DTOs =====

// CheckPermissionRequest represents the request to check user permissions
type CheckPermissionRequest struct {
	UserID         uint   `json:"user_id" binding:"required" example:"1"`
	Permission     string `json:"permission" binding:"required" example:"users.create"`
	Resource       string `json:"resource,omitempty" example:"users"`
	OrganizationID *uint  `json:"organization_id,omitempty" example:"1"`
	TeamID         *uint  `json:"team_id,omitempty" example:"1"`
}

// CheckPermissionResponse represents the permission check response
type CheckPermissionResponse struct {
	HasPermission bool     `json:"has_permission"`
	UserID        uint     `json:"user_id"`
	Permission    string   `json:"permission"`
	Resource      string   `json:"resource,omitempty"`
	Roles         []string `json:"roles"`
	Source        string   `json:"source"` // "global", "organization", "team"
}

// ===== User Permissions Summary DTOs =====

// UserPermissionsSummaryResponse represents a summary of user's permissions
type UserPermissionsSummaryResponse struct {
	UserID               uint                       `json:"user_id"`
	GlobalRoles          []RoleResponse             `json:"global_roles"`
	OrganizationRoles    []OrganizationRoleResponse `json:"organization_roles"`
	TeamRoles            []TeamRoleResponse         `json:"team_roles"`
	AllPermissions       []string                   `json:"all_permissions"`
	EffectivePermissions []PermissionResponse       `json:"effective_permissions"`
}

// ===== List DTOs =====

// ListQuery represents common list query parameters
type ListQuery struct {
	Page     int    `form:"page,default=1" example:"1"`
	PageSize int    `form:"page_size,default=20" example:"20"`
	Search   string `form:"search" example:"admin"`
	Status   *int   `form:"status" example:"1"`
	OrderBy  string `form:"order_by,default=created_at" example:"created_at"`
	Order    string `form:"order,default=desc" example:"desc"`
}

// ListRolesQuery represents query parameters for listing roles
type ListRolesQuery struct {
	ListQuery
	Level    *int  `form:"level" example:"100"`
	IsSystem *bool `form:"is_system" example:"false"`
}

// ListPermissionsQuery represents query parameters for listing permissions
type ListPermissionsQuery struct {
	ListQuery
	Resource string `form:"resource" example:"users"`
	Action   string `form:"action" example:"create"`
	Category string `form:"category" example:"user_management"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
