package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role with permissions
// @Tags roles
// @Accept json
// @Produce json
// @Param role body CreateRoleRequest true "Role data"
// @Success 201 {object} RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/roles [post]
func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create role object
	role := &Role{
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    req.Description,
		OrganizationID: req.OrganizationID,
		Permissions:    req.Permissions,
		IsDefault:      req.IsDefault,
	}

	if err := h.service.CreateRole(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := RoleResponse{
		ID:             role.ID,
		Name:           role.Name,
		DisplayName:    role.DisplayName,
		Description:    role.Description,
		OrganizationID: role.OrganizationID,
		Permissions:    role.Permissions,
		IsDefault:      role.IsDefault,
		CreatedAt:      role.CreatedAt,
		UpdatedAt:      role.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetRole godoc
// @Summary Get role by ID
// @Description Get details of a role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} RoleResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/roles/{id} [get]
func (h *Handler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	role, err := h.service.GetRole(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	// Convert to response format
	response := RoleResponse{
		ID:             role.ID,
		Name:           role.Name,
		DisplayName:    role.DisplayName,
		Description:    role.Description,
		OrganizationID: role.OrganizationID,
		Permissions:    role.Permissions,
		IsDefault:      role.IsDefault,
		CreatedAt:      role.CreatedAt,
		UpdatedAt:      role.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update a role's details and permissions
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body UpdateRoleRequest true "Role data"
// @Success 200 {object} RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/roles/{id} [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing role
	role, err := h.service.GetRole(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	// Update fields
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Permissions != "" {
		role.Permissions = req.Permissions
	}
	if req.IsDefault != nil {
		role.IsDefault = *req.IsDefault
	}

	if err := h.service.UpdateRole(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := RoleResponse{
		ID:             role.ID,
		Name:           role.Name,
		DisplayName:    role.DisplayName,
		Description:    role.Description,
		OrganizationID: role.OrganizationID,
		Permissions:    role.Permissions,
		IsDefault:      role.IsDefault,
		CreatedAt:      role.CreatedAt,
		UpdatedAt:      role.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/roles/{id} [delete]
func (h *Handler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListRoles godoc
// @Summary List roles in an organization
// @Description Get a list of roles in an organization
// @Tags roles
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Success 200 {array} RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/{organization_id}/roles [get]
func (h *Handler) ListRoles(c *gin.Context) {
	orgIDStr := c.Param("organization_id")
	orgIDVal, err := strconv.ParseUint(orgIDStr, 10, 32)
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
	
	roles, total, err := h.service.ListRoles(c.Request.Context(), uint(orgIDVal), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to response format
	var responses []RoleResponse
	for _, role := range roles {
		responses = append(responses, RoleResponse{
			ID:             role.ID,
			Name:           role.Name,
			DisplayName:    role.DisplayName,
			Description:    role.Description,
			OrganizationID: role.OrganizationID,
			Permissions:    role.Permissions,
			IsDefault:      role.IsDefault,
			CreatedAt:      role.CreatedAt,
			UpdatedAt:      role.UpdatedAt,
		})
	}
	
	c.JSON(http.StatusOK, PaginationResponse{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  responses,
	})
}
