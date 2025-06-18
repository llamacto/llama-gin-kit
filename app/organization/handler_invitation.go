package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateInvitation godoc
// @Summary Create a new invitation
// @Description Create a new invitation to join an organization
// @Tags invitations
// @Accept json
// @Produce json
// @Param invitation body CreateInvitationRequest true "Invitation data"
// @Success 201 {object} InvitationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invitations [post]
func (h *Handler) CreateInvitation(c *gin.Context) {
	var req CreateInvitationRequest
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

	// Create invitation object
	invitation := &Invitation{
		Email:          req.Email,
		OrganizationID: req.OrganizationID,
		TeamID:         req.TeamID,
		RoleID:         req.RoleID,
		InvitedBy:      invitedBy.(uint),
		Status:         0, // Pending
	}

	if err := h.service.InviteMember(c.Request.Context(), invitation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := InvitationResponse{
		ID:             invitation.ID,
		Email:          invitation.Email,
		OrganizationID: invitation.OrganizationID,
		TeamID:         invitation.TeamID,
		RoleID:         invitation.RoleID,
		InvitedBy:      invitation.InvitedBy,
		Token:          invitation.Token,
		ExpiresAt:      invitation.ExpiresAt,
		Status:         invitation.Status,
		CreatedAt:      invitation.CreatedAt,
		UpdatedAt:      invitation.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetInvitation godoc
// @Summary Get invitation by ID
// @Description Get details of an invitation by ID
// @Tags invitations
// @Accept json
// @Produce json
// @Param id path int true "Invitation ID"
// @Success 200 {object} InvitationResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invitations/{id} [get]
func (h *Handler) GetInvitation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	invitation, err := h.service.GetInvitation(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
		return
	}

	// Convert to response format
	response := InvitationResponse{
		ID:             invitation.ID,
		Email:          invitation.Email,
		OrganizationID: invitation.OrganizationID,
		TeamID:         invitation.TeamID,
		RoleID:         invitation.RoleID,
		InvitedBy:      invitation.InvitedBy,
		ExpiresAt:      invitation.ExpiresAt,
		Status:         invitation.Status,
		CreatedAt:      invitation.CreatedAt,
		UpdatedAt:      invitation.UpdatedAt,
	}

	// Don't include token in response for security
	c.JSON(http.StatusOK, response)
}

// CancelInvitation godoc
// @Summary Cancel an invitation
// @Description Cancel a pending invitation
// @Tags invitations
// @Accept json
// @Produce json
// @Param id path int true "Invitation ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invitations/{id} [delete]
func (h *Handler) CancelInvitation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.service.CancelInvitation(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AcceptInvitation godoc
// @Summary Accept an invitation
// @Description Accept an invitation to join an organization
// @Tags invitations
// @Accept json
// @Produce json
// @Param request body AcceptInvitationRequest true "Accept invitation request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invitations/accept [post]
func (h *Handler) AcceptInvitation(c *gin.Context) {
	var req AcceptInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.ProcessInvitation(c.Request.Context(), req.Token, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invitation accepted successfully"})
}

// ListInvitations godoc
// @Summary List invitations in an organization
// @Description Get a paginated list of invitations in an organization
// @Tags invitations
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/{organization_id}/invitations [get]
func (h *Handler) ListInvitations(c *gin.Context) {
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
	
	// Get invitations
	invitations, total, err := h.service.ListInvitations(c.Request.Context(), uint(orgID), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to response format
	var responses []InvitationResponse
	for _, invitation := range invitations {
		// Don't include token in list responses for security
		responses = append(responses, InvitationResponse{
			ID:             invitation.ID,
			Email:          invitation.Email,
			OrganizationID: invitation.OrganizationID,
			TeamID:         invitation.TeamID,
			RoleID:         invitation.RoleID,
			InvitedBy:      invitation.InvitedBy,
			ExpiresAt:      invitation.ExpiresAt,
			Status:         invitation.Status,
			CreatedAt:      invitation.CreatedAt,
			UpdatedAt:      invitation.UpdatedAt,
		})
	}
	
	c.JSON(http.StatusOK, PaginationResponse{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  responses,
	})
}

// GetInvitationByToken godoc
// @Summary Get invitation by token
// @Description Get details of an invitation by its token
// @Tags invitations
// @Accept json
// @Produce json
// @Param token path string true "Invitation Token"
// @Success 200 {object} InvitationResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invitations/token/{token} [get]
func (h *Handler) GetInvitationByToken(c *gin.Context) {
	token := c.Param("token")

	invitation, err := h.service.GetInvitationByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
		return
	}

	// Convert to response format (exclude token for security)
	response := InvitationResponse{
		ID:             invitation.ID,
		Email:          invitation.Email,
		OrganizationID: invitation.OrganizationID,
		TeamID:         invitation.TeamID,
		RoleID:         invitation.RoleID,
		InvitedBy:      invitation.InvitedBy,
		ExpiresAt:      invitation.ExpiresAt,
		Status:         invitation.Status,
		CreatedAt:      invitation.CreatedAt,
		UpdatedAt:      invitation.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}
