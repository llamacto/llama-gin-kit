package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new team within an organization
// @Tags teams
// @Accept json
// @Produce json
// @Param team body CreateTeamRequest true "Team data"
// @Success 201 {object} TeamResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/teams [post]
func (h *Handler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create team object
	team := &Team{
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    req.Description,
		OrganizationID: req.OrganizationID,
		ParentTeamID:   req.ParentTeamID,
		Settings:       req.Settings,
		Status:         1, // Active
	}

	if err := h.service.CreateTeam(c.Request.Context(), team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := TeamResponse{
		ID:             team.ID,
		Name:           team.Name,
		DisplayName:    team.DisplayName,
		Description:    team.Description,
		OrganizationID: team.OrganizationID,
		ParentTeamID:   team.ParentTeamID,
		Settings:       team.Settings,
		Status:         team.Status,
		CreatedAt:      team.CreatedAt,
		UpdatedAt:      team.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetTeam godoc
// @Summary Get team by ID
// @Description Get details of a team by its ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} TeamResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/teams/{id} [get]
func (h *Handler) GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	team, err := h.service.GetTeam(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	// Convert to response format
	response := TeamResponse{
		ID:             team.ID,
		Name:           team.Name,
		DisplayName:    team.DisplayName,
		Description:    team.Description,
		OrganizationID: team.OrganizationID,
		ParentTeamID:   team.ParentTeamID,
		Settings:       team.Settings,
		Status:         team.Status,
		CreatedAt:      team.CreatedAt,
		UpdatedAt:      team.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTeam godoc
// @Summary Update a team
// @Description Update a team's details
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param team body UpdateTeamRequest true "Team data"
// @Success 200 {object} TeamResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/teams/{id} [put]
func (h *Handler) UpdateTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing team
	team, err := h.service.GetTeam(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	// Update fields
	if req.DisplayName != "" {
		team.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		team.Description = req.Description
	}
	if req.ParentTeamID != nil {
		team.ParentTeamID = req.ParentTeamID
	}
	if req.Settings != "" {
		team.Settings = req.Settings
	}
	if req.Status != nil {
		team.Status = *req.Status
	}

	if err := h.service.UpdateTeam(c.Request.Context(), team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	response := TeamResponse{
		ID:             team.ID,
		Name:           team.Name,
		DisplayName:    team.DisplayName,
		Description:    team.Description,
		OrganizationID: team.OrganizationID,
		ParentTeamID:   team.ParentTeamID,
		Settings:       team.Settings,
		Status:         team.Status,
		CreatedAt:      team.CreatedAt,
		UpdatedAt:      team.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTeam godoc
// @Summary Delete a team
// @Description Delete a team by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/teams/{id} [delete]
func (h *Handler) DeleteTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.service.DeleteTeam(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListTeams godoc
// @Summary List teams in an organization
// @Description Get a paginated list of teams in an organization
// @Tags teams
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/organizations/{organization_id}/teams [get]
func (h *Handler) ListTeams(c *gin.Context) {
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
	
	// Get teams
	teams, total, err := h.service.ListTeams(c.Request.Context(), uint(orgID), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to response format
	var responses []TeamResponse
	for _, team := range teams {
		responses = append(responses, TeamResponse{
			ID:             team.ID,
			Name:           team.Name,
			DisplayName:    team.DisplayName,
			Description:    team.Description,
			OrganizationID: team.OrganizationID,
			ParentTeamID:   team.ParentTeamID,
			Settings:       team.Settings,
			Status:         team.Status,
			CreatedAt:      team.CreatedAt,
			UpdatedAt:      team.UpdatedAt,
		})
	}
	
	c.JSON(http.StatusOK, PaginationResponse{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  responses,
	})
}
