package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler for organization endpoints
type Handler struct {
	service OrganizationService
}

// NewHandler creates a new organization handler
func NewHandler(service OrganizationService) *Handler {
	return &Handler{service: service}
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Create a new organization and set creator as admin
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body CreateOrganizationRequest true "Organization data"
// @Success 201 {object} OrganizationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations [post]
func (h *Handler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
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

	// Create organization object
	org := &Organization{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Logo:        req.Logo,
		Website:     req.Website,
		Settings:    req.Settings,
		Status:      1, // Active
	}

	if err := h.service.CreateOrganization(c.Request.Context(), org, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		DisplayName: org.DisplayName,
		Description: org.Description,
		Logo:        org.Logo,
		Website:     org.Website,
		Settings:    org.Settings,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetOrganization godoc
// @Summary Get organization by ID
// @Description Get details of an organization by its ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} OrganizationResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/{id} [get]
func (h *Handler) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	org, err := h.service.GetOrganization(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}

	// Convert to response format
	response := OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		DisplayName: org.DisplayName,
		Description: org.Description,
		Logo:        org.Logo,
		Website:     org.Website,
		Settings:    org.Settings,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateOrganization godoc
// @Summary Update an organization
// @Description Update an organization's details
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param organization body UpdateOrganizationRequest true "Organization data"
// @Success 200 {object} OrganizationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/{id} [put]
func (h *Handler) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing organization
	org, err := h.service.GetOrganization(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}

	// Update fields
	if req.DisplayName != "" {
		org.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		org.Description = req.Description
	}
	if req.Logo != "" {
		org.Logo = req.Logo
	}
	if req.Website != "" {
		org.Website = req.Website
	}
	if req.Settings != "" {
		org.Settings = req.Settings
	}
	if req.Status != nil {
		org.Status = *req.Status
	}

	if err := h.service.UpdateOrganization(c.Request.Context(), org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		DisplayName: org.DisplayName,
		Description: org.Description,
		Logo:        org.Logo,
		Website:     org.Website,
		Settings:    org.Settings,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteOrganization godoc
// @Summary Delete an organization
// @Description Delete an organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/{id} [delete]
func (h *Handler) DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.service.DeleteOrganization(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListOrganizations godoc
// @Summary List organizations
// @Description Get a paginated list of organizations
// @Tags organizations
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} PaginationResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations [get]
func (h *Handler) ListOrganizations(c *gin.Context) {
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
	
	// Get organizations
	orgs, total, err := h.service.ListOrganizations(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to response format
	var responses []OrganizationResponse
	for _, org := range orgs {
		responses = append(responses, OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			DisplayName: org.DisplayName,
			Description: org.Description,
			Logo:        org.Logo,
			Website:     org.Website,
			Settings:    org.Settings,
			Status:      org.Status,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		})
	}
	
	c.JSON(http.StatusOK, PaginationResponse{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  responses,
	})
}

// GetMyOrganizations godoc
// @Summary List user's organizations
// @Description Get organizations the current user is a member of
// @Tags organizations
// @Accept json
// @Produce json
// @Success 200 {array} OrganizationResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/me [get]
func (h *Handler) GetMyOrganizations(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	orgs, err := h.service.GetUserOrganizations(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to response format
	var responses []OrganizationResponse
	for _, org := range orgs {
		responses = append(responses, OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			DisplayName: org.DisplayName,
			Description: org.Description,
			Logo:        org.Logo,
			Website:     org.Website,
			Settings:    org.Settings,
			Status:      org.Status,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		})
	}
	
	c.JSON(http.StatusOK, responses)
}

// CheckPermission godoc
// @Summary Check user permission
// @Description Check if current user has a specific permission in an organization
// @Tags permissions
// @Accept json
// @Produce json
// @Param request body CheckPermissionRequest true "Permission check request"
// @Success 200 {object} CheckPermissionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/permissions/check [post]
func (h *Handler) CheckPermission(c *gin.Context) {
	var req CheckPermissionRequest
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
	
	hasPermission, err := h.service.CheckPermission(
		c.Request.Context(), 
		userID.(uint), 
		req.OrganizationID, 
		req.Permission,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, CheckPermissionResponse{HasPermission: hasPermission})
}
