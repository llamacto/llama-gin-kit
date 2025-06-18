package team

// CreateTeamRequest represents the request payload for creating a team
type CreateTeamRequest struct {
	Name           string `json:"name" binding:"required,min=2,max=100"`
	DisplayName    string `json:"display_name" binding:"max=100"`
	Description    string `json:"description" binding:"max=500"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	ParentTeamID   *uint  `json:"parent_team_id"`
	// Settings       string `json:"settings"` // Temporarily disabled
}

// UpdateTeamRequest represents the request payload for updating a team
type UpdateTeamRequest struct {
	Name         string `json:"name" binding:"min=2,max=100"`
	DisplayName  string `json:"display_name" binding:"max=100"`
	Description  string `json:"description" binding:"max=500"`
	ParentTeamID *uint  `json:"parent_team_id"`
	// Settings     string `json:"settings"` // Temporarily disabled
	Status *int `json:"status"`
}

// TeamResponse represents the response structure for team data
type TeamResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
	OrganizationID uint   `json:"organization_id"`
	ParentTeamID   *uint  `json:"parent_team_id"`
	// Settings       string `json:"settings"` // Temporarily disabled
	Status      int    `json:"status"`
	MemberCount int64  `json:"member_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// TeamListResponse represents the response structure for team list
type TeamListResponse struct {
	Teams      []TeamResponse `json:"teams"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// TeamHierarchyResponse represents the response structure for team hierarchy
type TeamHierarchyResponse struct {
	Team     TeamResponse   `json:"team"`
	Parent   *TeamResponse  `json:"parent,omitempty"`
	Children []TeamResponse `json:"children,omitempty"`
}
