package role

// CreateRoleRequest represents the request payload for creating a role
type CreateRoleRequest struct {
	Name           string   `json:"name" binding:"required,min=2,max=50"`
	DisplayName    string   `json:"display_name" binding:"max=100"`
	Description    string   `json:"description" binding:"max=255"`
	OrganizationID *uint    `json:"organization_id"`
	Permissions    []string `json:"permissions"`
	IsDefault      bool     `json:"is_default"`
}

// UpdateRoleRequest represents the request payload for updating a role
type UpdateRoleRequest struct {
	Name        string   `json:"name" binding:"min=2,max=50"`
	DisplayName string   `json:"display_name" binding:"max=100"`
	Description string   `json:"description" binding:"max=255"`
	Permissions []string `json:"permissions"`
	IsDefault   *bool    `json:"is_default"`
}

// RoleResponse represents the response structure for role data
type RoleResponse struct {
	ID             uint     `json:"id"`
	Name           string   `json:"name"`
	DisplayName    string   `json:"display_name"`
	Description    string   `json:"description"`
	OrganizationID *uint    `json:"organization_id"`
	Permissions    []string `json:"permissions"`
	IsDefault      bool     `json:"is_default"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// RoleListResponse represents the response structure for role list
type RoleListResponse struct {
	Roles      []RoleResponse `json:"roles"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// PermissionResponse represents the response structure for permission data
type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// PermissionListResponse represents the response structure for permission list
type PermissionListResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
	Total       int64                `json:"total"`
}

// PermissionsByCategoryResponse represents permissions grouped by category
type PermissionsByCategoryResponse struct {
	Categories []PermissionCategoryResponse `json:"categories"`
}

// PermissionCategoryResponse represents a category with its permissions
type PermissionCategoryResponse struct {
	Category    string               `json:"category"`
	Permissions []PermissionResponse `json:"permissions"`
}
