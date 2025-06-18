package apikey

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/pkg/response"
)

// Handler interface for API key operations
type Handler interface {
	// Create creates a new API key
	Create(c *gin.Context)
	
	// Get gets an API key by ID
	Get(c *gin.Context)
	
	// List lists all API keys for the authenticated user
	List(c *gin.Context)
	
	// Update updates an API key
	Update(c *gin.Context)
	
	// Delete revokes (deletes) an API key
	Delete(c *gin.Context)
}

// handler implements the Handler interface
type handler struct {
	service Service
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(service Service) Handler {
	return &handler{service: service}
}

// Create creates a new API key
// @Summary Create a new API key
// @Description Creates a new API key for the authenticated user
// @Tags API Keys
// @Accept json
// @Produce json
// @Param request body CreateRequest true "API Key Details"
// @Success 201 {object} Response "API Key created"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/apikeys [post]
// @Security BearerAuth
func (h *handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Set expiry time
	var expiry *time.Time
	if !req.NeverExpire {
		if !req.ExpiresAt.IsZero() {
			expiry = &req.ExpiresAt
		} else {
			// Default expiry: 1 year
			defaultExpiry := time.Now().AddDate(1, 0, 0)
			expiry = &defaultExpiry
		}
	}

	// Generate API key
	key, apiKey, err := h.service.GenerateAPIKey(userID.(uint), req.Name, expiry, req.Permissions)
	if err != nil {
		response.InternalServerError(c, "Failed to create API key", err)
		return
	}

	// Convert to response DTO
	resp := ToResponse(apiKey, key)

	// Return response
	c.JSON(http.StatusCreated, resp)
}

// Get gets an API key by ID
// @Summary Get an API key
// @Description Gets an API key by its ID
// @Tags API Keys
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 200 {object} Response "API Key details"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/apikeys/{id} [get]
// @Security BearerAuth
func (h *handler) Get(c *gin.Context) {
	// Parse API key ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid API key ID", err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get API key
	apiKey, err := h.service.GetAPIKey(uint(id))
	if err != nil {
		response.NotFound(c, "API key not found", err)
		return
	}

	// Security check: ensure the key belongs to the user
	if apiKey.UserID != userID.(uint) {
		response.Unauthorized(c, "You do not have permission to access this API key")
		return
	}

	// Convert to response DTO
	resp := ToResponse(apiKey, "")

	// Return response
	c.JSON(http.StatusOK, resp)
}

// List lists all API keys for the authenticated user
// @Summary List API keys
// @Description Lists all API keys for the authenticated user with pagination
// @Tags API Keys
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} ListResponse "List of API keys"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/apikeys [get]
// @Security BearerAuth
func (h *handler) List(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get API keys
	apiKeys, total, err := h.service.ListAPIKeys(userID.(uint), page, perPage)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve API keys", err)
		return
	}

	// Convert to response DTO
	resp := ListResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    ToResponseList(apiKeys),
	}

	// Return response
	c.JSON(http.StatusOK, resp)
}

// Update updates an API key
// @Summary Update an API key
// @Description Updates an API key's name, permissions or expiry
// @Tags API Keys
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Param request body UpdateRequest true "API Key Details"
// @Success 200 {object} Response "Updated API key"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/apikeys/{id} [put]
// @Security BearerAuth
func (h *handler) Update(c *gin.Context) {
	// Parse API key ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid API key ID", err)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Set expiry time
	var expiry *time.Time
	if req.NeverExpire {
		expiry = nil
	} else if !req.ExpiresAt.IsZero() {
		expiry = &req.ExpiresAt
	}

	// Update API key
	apiKey, err := h.service.UpdateAPIKey(uint(id), userID.(uint), req.Name, expiry, req.Permissions)
	if err != nil {
		response.HandleError(c, "Failed to update API key", err)
		return
	}

	// Convert to response DTO
	resp := ToResponse(apiKey, "")

	// Return response
	c.JSON(http.StatusOK, resp)
}

// Delete revokes (deletes) an API key
// @Summary Delete an API key
// @Description Revokes (deletes) an API key
// @Tags API Keys
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 204 "No content"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/apikeys/{id} [delete]
// @Security BearerAuth
func (h *handler) Delete(c *gin.Context) {
	// Parse API key ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid API key ID", err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Delete API key
	if err := h.service.RevokeAPIKey(uint(id), userID.(uint)); err != nil {
		response.HandleError(c, "Failed to delete API key", err)
		return
	}

	// Return response
	c.Status(http.StatusNoContent)
}
