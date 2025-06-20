package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/pkg/response"
)

// Middleware provides a middleware for checking permissions.
type Middleware struct {
	service Service
}

// NewMiddleware creates a new authorization middleware.
func NewMiddleware(service Service) *Middleware {
	return &Middleware{service: service}
}

// RequirePermission creates a Gin middleware that checks if the user has a specific permission.
func (m *Middleware) RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		hasPerm, err := m.service.HasPermission(userID, requiredPermission)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to check permission")
			c.Abort()
			return
		}

		if !hasPerm {
			// Check for organization or team specific permissions if context is available
			// For simplicity, this basic check only covers global permissions.
			// A more advanced version could extract org/team IDs from the context or path.

			// Let's try to check for super_admin role as an override
			roles, err := m.service.GetUserRoles(userID)
			if err == nil {
				for _, role := range roles {
					if role.Role.Name == "super_admin" {
						c.Next()
						return
					}
				}
			}

			response.Error(c, http.StatusForbidden, "You do not have permission to perform this action")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates a Gin middleware that checks if the user has a specific role.
func (m *Middleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		userRoles, err := m.service.GetUserRoles(userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to check user roles")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range userRoles {
			if role.Role.Name == requiredRole || role.Role.Name == "super_admin" {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Error(c, http.StatusForbidden, "You do not have the required role to perform this action")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireLevel creates a Gin middleware that checks if the user's highest role level is sufficient.
func (m *Middleware) RequireLevel(requiredLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		userRoles, err := m.service.GetUserRoles(userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to check user roles")
			c.Abort()
			return
		}

		maxLevel := 0
		for _, role := range userRoles {
			if role.Role.Level > maxLevel {
				maxLevel = role.Role.Level
			}
		}

		if maxLevel < requiredLevel {
			response.Error(c, http.StatusForbidden, "Your role level is not high enough for this action")
			c.Abort()
			return
		}

		c.Next()
	}
}
