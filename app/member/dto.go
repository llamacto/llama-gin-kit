package member

// AddMemberRequest represents the request payload for adding a member to organization/team
type AddMemberRequest struct {
	UserID         uint  `json:"user_id" binding:"required"`
	OrganizationID uint  `json:"organization_id" binding:"required"`
	TeamID         *uint `json:"team_id"`
	RoleID         uint  `json:"role_id" binding:"required"`
}

// UpdateMemberRequest represents the request payload for updating member info
type UpdateMemberRequest struct {
	TeamID *uint `json:"team_id"`
	RoleID *uint `json:"role_id"`
	Status *int  `json:"status"`
}

// MemberResponse represents the response structure for member data
type MemberResponse struct {
	ID               uint   `json:"id"`
	UserID           uint   `json:"user_id"`
	UserName         string `json:"user_name"`
	UserEmail        string `json:"user_email"`
	UserNickname     string `json:"user_nickname"`
	UserAvatar       string `json:"user_avatar"`
	OrganizationID   uint   `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	TeamID           *uint  `json:"team_id"`
	TeamName         string `json:"team_name"`
	RoleID           uint   `json:"role_id"`
	RoleName         string `json:"role_name"`
	RoleDisplayName  string `json:"role_display_name"`
	Status           int    `json:"status"`
	JoinedAt         string `json:"joined_at"`
	InvitedBy        uint   `json:"invited_by"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// MemberListResponse represents the response structure for member list
type MemberListResponse struct {
	Members    []MemberResponse `json:"members"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// MemberStatsResponse represents the response structure for member statistics
type MemberStatsResponse struct {
	TotalMembers    int64 `json:"total_members"`
	ActiveMembers   int64 `json:"active_members"`
	PendingInvites  int64 `json:"pending_invites"`
	DisabledMembers int64 `json:"disabled_members"`
}
