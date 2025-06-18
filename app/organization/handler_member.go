package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddMember godoc
// @Summary Add a member to an organization
// @Description Add a new member to an organization with a specific role
// @Tags members
// @Accept json
// @Produce json
// @Param member body AddMemberRequest true "Member data"
// @Success 201 {object} MemberResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/members [post]
func (h *Handler) AddMember(c *gin.Context) {
	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware) for invite tracking
	invitedBy, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create member object
	member := &Member{
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		TeamID:         req.TeamID,
		RoleID:         req.RoleID,
		Status:         1, // Active
		InvitedBy:      invitedBy.(uint),
	}

	if err := h.service.AddMember(c.Request.Context(), member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := MemberResponse{
		ID:             member.ID,
		UserID:         member.UserID,
		OrganizationID: member.OrganizationID,
		TeamID:         member.TeamID,
		RoleID:         member.RoleID,
		Status:         member.Status,
		JoinedAt:       member.JoinedAt,
		InvitedBy:      member.InvitedBy,
		CreatedAt:      member.CreatedAt,
		UpdatedAt:      member.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetMember godoc
// @Summary Get member by ID
// @Description Get details of a member by ID
// @Tags members
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {object} MemberResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/members/{id} [get]
func (h *Handler) GetMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	member, err := h.service.GetMember(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	// Convert to response format
	response := MemberResponse{
		ID:             member.ID,
		UserID:         member.UserID,
		OrganizationID: member.OrganizationID,
		TeamID:         member.TeamID,
		RoleID:         member.RoleID,
		Status:         member.Status,
		JoinedAt:       member.JoinedAt,
		InvitedBy:      member.InvitedBy,
		CreatedAt:      member.CreatedAt,
		UpdatedAt:      member.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateMember godoc
// @Summary Update a member
// @Description Update a member's role or team
// @Tags members
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Param member body UpdateMemberRequest true "Member data"
// @Success 200 {object} MemberResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/members/{id} [put]
func (h *Handler) UpdateMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing member
	member, err := h.service.GetMember(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	// Update fields
	if req.TeamID != nil {
		member.TeamID = req.TeamID
	}
	member.RoleID = req.RoleID
	if req.Status != nil {
		member.Status = *req.Status
	}

	if err := h.service.UpdateMember(c.Request.Context(), member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := MemberResponse{
		ID:             member.ID,
		UserID:         member.UserID,
		OrganizationID: member.OrganizationID,
		TeamID:         member.TeamID,
		RoleID:         member.RoleID,
		Status:         member.Status,
		JoinedAt:       member.JoinedAt,
		InvitedBy:      member.InvitedBy,
		CreatedAt:      member.CreatedAt,
		UpdatedAt:      member.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// RemoveMember godoc
// @Summary Remove a member
// @Description Remove a member from an organization
// @Tags members
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/members/{id} [delete]
func (h *Handler) RemoveMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListMembers godoc
// @Summary List members in an organization
// @Description Get a paginated list of members in an organization
// @Tags members
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param team_id query int false "Filter by team ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/{organization_id}/members [get]
func (h *Handler) ListMembers(c *gin.Context) {
	orgIDStr := c.Param("organization_id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID format"})
		return
	}

	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}
	
	// Get members
	members, total, err := h.service.ListMembers(c.Request.Context(), uint(orgID), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to response format
	var responses []MemberResponse
	for _, member := range members {
		responses = append(responses, MemberResponse{
			ID:             member.ID,
			UserID:         member.UserID,
			OrganizationID: member.OrganizationID,
			TeamID:         member.TeamID,
			RoleID:         member.RoleID,
			Status:         member.Status,
			JoinedAt:       member.JoinedAt,
			InvitedBy:      member.InvitedBy,
			CreatedAt:      member.CreatedAt,
			UpdatedAt:      member.UpdatedAt,
		})
	}
	
	c.JSON(http.StatusOK, PaginationResponse{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  responses,
	})
}
