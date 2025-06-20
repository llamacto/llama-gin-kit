package authorization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/pkg/response"
)

// Handler for authorization endpoints
type Handler struct {
	service Service
}

// NewHandler creates a new authorization handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ===== Role Handlers =====

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role in the system.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param role body CreateRoleRequest true "Role details"
// @Success 201 {object} utils.SuccessResponse{data=RoleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles [post]
// @Security ApiKeyAuth
func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	role, err := h.service.CreateRole(req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, role)
}

// GetRole godoc
// @Summary Get a role by ID
// @Description Get details of a specific role by its ID.
// @Tags Authorization
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse{data=RoleResponse}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) GetRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.service.GetRole(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Role not found")
		return
	}

	response.Success(c, role)
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update an existing role's details.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body UpdateRoleRequest true "Role details to update"
// @Success 200 {object} utils.SuccessResponse{data=RoleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles/{id} [put]
// @Security ApiKeyAuth
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	role, err := h.service.UpdateRole(uint(id), req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, role)
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by its ID.
// @Tags Authorization
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.service.DeleteRole(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListRoles godoc
// @Summary List roles
// @Description Get a paginated list of roles.
// @Tags Authorization
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search query"
// @Param status query int false "Role status"
// @Success 200 {object} utils.SuccessResponse{data=ListResponse}
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles [get]
// @Security ApiKeyAuth
func (h *Handler) ListRoles(c *gin.Context) {
	var query ListRolesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	roles, err := h.service.ListRoles(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, roles)
}

// GetRoleWithPermissions godoc
// @Summary Get a role with its permissions
// @Description Get a role and all associated permissions.
// @Tags Authorization
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse{data=RoleWithPermissionsResponse}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles/{id}/permissions [get]
// @Security ApiKeyAuth
func (h *Handler) GetRoleWithPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	roleWithPerms, err := h.service.GetRoleWithPermissions(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, roleWithPerms)
}

// ===== Permission Handlers =====

// CreatePermission godoc
// @Summary Create a new permission
// @Description Create a new permission in the system.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param permission body CreatePermissionRequest true "Permission details"
// @Success 201 {object} utils.SuccessResponse{data=PermissionResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/permissions [post]
// @Security ApiKeyAuth
func (h *Handler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	permission, err := h.service.CreatePermission(req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, permission)
}

// ListPermissions godoc
// @Summary List permissions
// @Description Get a paginated list of permissions.
// @Tags Authorization
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search query"
// @Param resource query string false "Resource name"
// @Param action query string false "Action name"
// @Param category query string false "Category name"
// @Success 200 {object} utils.SuccessResponse{data=ListResponse}
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/permissions [get]
// @Security ApiKeyAuth
func (h *Handler) ListPermissions(c *gin.Context) {
	var query ListPermissionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	permissions, err := h.service.ListPermissions(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, permissions)
}

// ===== Assignment Handlers =====

// AssignPermissionsToRole godoc
// @Summary Assign permissions to a role
// @Description Assign a list of permissions to a role.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param permissions body AssignPermissionsRequest true "Permission IDs to assign"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles/{id}/permissions [post]
// @Security ApiKeyAuth
func (h *Handler) AssignPermissionsToRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.service.AssignPermissionsToRole(uint(roleID), req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// RemovePermissionsFromRole godoc
// @Summary Remove permissions from a role
// @Description Remove a list of permissions from a role.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param permissions body RemovePermissionsRequest true "Permission IDs to remove"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/roles/{id}/permissions [delete]
// @Security ApiKeyAuth
func (h *Handler) RemovePermissionsFromRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req RemovePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.service.RemovePermissionsFromRole(uint(roleID), req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// AssignRoleToUser godoc
// @Summary Assign a role to a user
// @Description Assign a role to a specific user.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param assignment body AssignRoleRequest true "User and Role IDs"
// @Success 200 {object} utils.SuccessResponse{data=UserRoleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/users/roles [post]
// @Security ApiKeyAuth
func (h *Handler) AssignRoleToUser(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	userRole, err := h.service.AssignRoleToUser(req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, userRole)
}

// RemoveRoleFromUser godoc
// @Summary Remove a role from a user
// @Description Remove a role from a specific user.
// @Tags Authorization
// @Produce json
// @Param userId path int true "User ID"
// @Param roleId path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/users/{userId}/roles/{roleId} [delete]
// @Security ApiKeyAuth
func (h *Handler) RemoveRoleFromUser(c *gin.Context) {
	userIDParam, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roleID, err := strconv.Atoi(c.Param("roleId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	removedBy, err := getUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.service.RemoveRoleFromUser(uint(userIDParam), uint(roleID), removedBy)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUserRoles godoc
// @Summary Get a user's roles
// @Description Get a list of roles assigned to a user.
// @Tags Authorization
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} utils.SuccessResponse{data=[]UserRoleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/users/{userId}/roles [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roles, err := h.service.GetUserRoles(uint(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, roles)
}

// CheckPermission godoc
// @Summary Check user permission
// @Description Check if a user has a specific permission.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param check body CheckPermissionRequest true "Permission check details"
// @Success 200 {object} utils.SuccessResponse{data=CheckPermissionResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/check-permission [post]
// @Security ApiKeyAuth
func (h *Handler) CheckPermission(c *gin.Context) {
	var req CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.service.CheckPermission(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetUserPermissionsSummary godoc
// @Summary Get user permissions summary
// @Description Get a summary of all permissions for a user.
// @Tags Authorization
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} utils.SuccessResponse{data=UserPermissionsSummaryResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/users/{userId}/permissions-summary [get]
// @Security ApiKeyAuth
func (h *Handler) GetUserPermissionsSummary(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	summary, err := h.service.GetUserPermissionsSummary(uint(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, summary)
}
