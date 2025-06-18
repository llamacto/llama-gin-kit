package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler struct for organization operations
type Handler struct {
	service Service
}

// NewHandler creates a new organization handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreateOrganization creates a new organization without settings
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

	org := &Organization{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Logo:        req.Logo,
		Website:     req.Website,
		Status:      1, // Active
	}

	if err := h.service.CreateOrganization(c.Request.Context(), org, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format (without settings)
	response := gin.H{
		"id":           org.ID,
		"name":         org.Name,
		"display_name": org.DisplayName,
		"description":  org.Description,
		"logo":         org.Logo,
		"website":      org.Website,
		"status":       org.Status,
		"created_at":   org.CreatedAt,
		"updated_at":   org.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetOrganization gets an organization by ID
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

	response := gin.H{
		"id":           org.ID,
		"name":         org.Name,
		"display_name": org.DisplayName,
		"description":  org.Description,
		"logo":         org.Logo,
		"website":      org.Website,
		"status":       org.Status,
		"created_at":   org.CreatedAt,
		"updated_at":   org.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// ListOrganizations lists organizations with pagination
func (h *Handler) ListOrganizations(c *gin.Context) {
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

	orgs, total, err := h.service.ListOrganizations(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []gin.H
	for _, org := range orgs {
		responses = append(responses, gin.H{
			"id":           org.ID,
			"name":         org.Name,
			"display_name": org.DisplayName,
			"description":  org.Description,
			"logo":         org.Logo,
			"website":      org.Website,
			"status":       org.Status,
			"created_at":   org.CreatedAt,
			"updated_at":   org.UpdatedAt,
		})
	}

	response := gin.H{
		"total": total,
		"page":  page,
		"size":  size,
		"data":  responses,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateOrganization updates an organization
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
	if req.Status != nil {
		org.Status = *req.Status
	}

	if err := h.service.UpdateOrganization(c.Request.Context(), org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"id":           org.ID,
		"name":         org.Name,
		"display_name": org.DisplayName,
		"description":  org.Description,
		"logo":         org.Logo,
		"website":      org.Website,
		"status":       org.Status,
		"created_at":   org.CreatedAt,
		"updated_at":   org.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteOrganization deletes an organization
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

// GetMyOrganizations gets organizations for the current user
func (h *Handler) GetMyOrganizations(c *gin.Context) {
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

	var responses []gin.H
	for _, org := range orgs {
		responses = append(responses, gin.H{
			"id":           org.ID,
			"name":         org.Name,
			"display_name": org.DisplayName,
			"description":  org.Description,
			"logo":         org.Logo,
			"website":      org.Website,
			"status":       org.Status,
			"created_at":   org.CreatedAt,
			"updated_at":   org.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}
