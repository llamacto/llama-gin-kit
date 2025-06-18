package invitation

// CreateInvitationRequest represents the request payload for creating an invitation
type CreateInvitationRequest struct {
	Email          string `json:"email" binding:"required,email"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	TeamID         *uint  `json:"team_id"`
	RoleID         uint   `json:"role_id" binding:"required"`
}

// BatchInvitationRequest represents the request payload for batch invitations
type BatchInvitationRequest struct {
	Emails         []string `json:"emails" binding:"required,min=1"`
	OrganizationID uint     `json:"organization_id" binding:"required"`
	TeamID         *uint    `json:"team_id"`
	RoleID         uint     `json:"role_id" binding:"required"`
}

// AcceptInvitationRequest represents the request payload for accepting an invitation
type AcceptInvitationRequest struct {
	Token string `json:"token" binding:"required"`
}

// ResendInvitationRequest represents the request payload for resending an invitation
type ResendInvitationRequest struct {
	InvitationID uint `json:"invitation_id" binding:"required"`
}

// InvitationResponse represents the response structure for invitation data
type InvitationResponse struct {
	ID               uint   `json:"id"`
	Email            string `json:"email"`
	OrganizationID   uint   `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	TeamID           *uint  `json:"team_id"`
	TeamName         string `json:"team_name"`
	RoleID           uint   `json:"role_id"`
	RoleName         string `json:"role_name"`
	RoleDisplayName  string `json:"role_display_name"`
	InvitedBy        uint   `json:"invited_by"`
	InviterName      string `json:"inviter_name"`
	InviterEmail     string `json:"inviter_email"`
	Token            string `json:"token"`
	ExpiresAt        string `json:"expires_at"`
	Status           int    `json:"status"`
	StatusText       string `json:"status_text"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// InvitationListResponse represents the response structure for invitation list
type InvitationListResponse struct {
	Invitations []InvitationResponse `json:"invitations"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	TotalPages  int                  `json:"total_pages"`
}

// InvitationStatsResponse represents the response structure for invitation statistics
type InvitationStatsResponse struct {
	Total    int64 `json:"total"`
	Pending  int64 `json:"pending"`
	Accepted int64 `json:"accepted"`
	Rejected int64 `json:"rejected"`
	Expired  int64 `json:"expired"`
}

// BatchInvitationResponse represents the response structure for batch invitation results
type BatchInvitationResponse struct {
	Success []InvitationResponse `json:"success"`
	Failed  []BatchFailedResult  `json:"failed"`
	Summary BatchSummary         `json:"summary"`
}

// BatchFailedResult represents a failed invitation in batch operation
type BatchFailedResult struct {
	Email  string `json:"email"`
	Reason string `json:"reason"`
}

// BatchSummary represents summary of batch invitation operation
type BatchSummary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}
