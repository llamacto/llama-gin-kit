package team

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/pkg/response"
)

// Handler defines the interface for team HTTP handlers
type Handler interface {
	CreateTeam(c *gin.Context)
	GetTeam(c *gin.Context)
	GetTeamsByOrganization(c *gin.Context)
	UpdateTeam(c *gin.Context)
	DeleteTeam(c *gin.Context)
	GetTeamHierarchy(c *gin.Context)
}

// handler implements the Handler interface
type handler struct {
	service Service
}

// NewHandler creates a new team handler instance
func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// CreateTeam creates a new team
// @Summary Create a new team
// @Description Create a new team within an organization
// @Tags teams
// @Accept json
// @Produce json
// @Param request body CreateTeamRequest true "Team creation request"
// @Success 201 {object} response.Response{data=TeamResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/teams [post]
func (h *handler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	team, err := h.service.CreateTeam(&req, userIDUint)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create team")
		return
	}

	response.Success(c, team)
}

// GetTeam retrieves a team by ID
// @Summary Get team by ID
// @Description Get team details by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} response.Response{data=TeamResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/teams/{id} [get]
func (h *handler) GetTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	team, err := h.service.GetTeamByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Team not found")
		return
	}

	response.Success(c, team)
}

// GetTeamsByOrganization retrieves teams by organization ID
// @Summary Get teams by organization
// @Description Get all teams within an organization with pagination
// @Tags teams
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} response.Response{data=TeamListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organization_id}/teams [get]
func (h *handler) GetTeamsByOrganization(c *gin.Context) {
	orgIDParam := c.Param("organization_id")
	organizationID, err := strconv.ParseUint(orgIDParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid organization ID")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	teams, err := h.service.GetTeamsByOrganization(uint(organizationID), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve teams")
		return
	}

	response.Success(c, teams)
}

// UpdateTeam updates a team
// @Summary Update team
// @Description Update team information
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param request body UpdateTeamRequest true "Team update request"
// @Success 200 {object} response.Response{data=TeamResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/teams/{id} [put]
func (h *handler) UpdateTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	team, err := h.service.UpdateTeam(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update team")
		return
	}

	response.Success(c, team)
}

// DeleteTeam deletes a team
// @Summary Delete team
// @Description Delete a team
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/teams/{id} [delete]
func (h *handler) DeleteTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	err = h.service.DeleteTeam(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete team")
		return
	}

	response.Success(c, nil)
}

// GetTeamHierarchy retrieves team hierarchy
// @Summary Get team hierarchy
// @Description Get team hierarchy with parent and children
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} response.Response{data=TeamHierarchyResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/teams/{id}/hierarchy [get]
func (h *handler) GetTeamHierarchy(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	hierarchy, err := h.service.GetTeamHierarchy(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Team hierarchy not found")
		return
	}

	response.Success(c, hierarchy)
}
