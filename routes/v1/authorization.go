package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/authorization"
	"github.com/llamacto/llama-gin-kit/pkg/middleware"
)

// RegisterAuthRoutes registers all authorization routes.
func RegisterAuthRoutes(router *gin.RouterGroup, authService authorization.Service) {
	// Initialize handler and middleware
	handler := authorization.NewHandler(authService)
	authMiddleware := authorization.NewMiddleware(authService)

	// Create a new group for authorization routes, e.g., /v1/auth
	authGroup := router.Group("/auth")
	authGroup.Use(middleware.JWTAuth()) // All routes in this group require authentication

	// Role management routes
	rolesGroup := authGroup.Group("/roles")
	rolesGroup.Use(authMiddleware.RequireRole("admin")) // Only admins can manage roles
	{
		rolesGroup.POST("", authMiddleware.RequirePermission("roles.create"), handler.CreateRole)
		rolesGroup.GET("", authMiddleware.RequirePermission("roles.read"), handler.ListRoles)
		rolesGroup.GET("/:id", authMiddleware.RequirePermission("roles.read"), handler.GetRole)
		rolesGroup.PUT("/:id", authMiddleware.RequirePermission("roles.update"), handler.UpdateRole)
		rolesGroup.DELETE("/:id", authMiddleware.RequirePermission("roles.delete"), handler.DeleteRole)

		// Role-Permission assignment routes
		rolesGroup.GET("/:id/permissions", authMiddleware.RequirePermission("roles.read"), handler.GetRoleWithPermissions)
		rolesGroup.POST("/:id/permissions", authMiddleware.RequirePermission("roles.assign_permissions"), handler.AssignPermissionsToRole)
		rolesGroup.DELETE("/:id/permissions", authMiddleware.RequirePermission("roles.remove_permissions"), handler.RemovePermissionsFromRole)
	}

	// Permission management routes
	permissionsGroup := authGroup.Group("/permissions")
	permissionsGroup.Use(authMiddleware.RequireRole("admin")) // Only admins can manage permissions
	{
		permissionsGroup.POST("", authMiddleware.RequirePermission("permissions.create"), handler.CreatePermission)
		permissionsGroup.GET("", authMiddleware.RequirePermission("permissions.read"), handler.ListPermissions)
	}

	// User-Role assignment routes
	usersGroup := authGroup.Group("/users")
	usersGroup.Use(authMiddleware.RequireRole("admin")) // Only admins can manage user roles
	{
		usersGroup.POST("/roles", authMiddleware.RequirePermission("users.assign_role"), handler.AssignRoleToUser)
		usersGroup.GET("/:userId/roles", authMiddleware.RequirePermission("users.read_roles"), handler.GetUserRoles)
		usersGroup.DELETE("/:userId/roles/:roleId", authMiddleware.RequirePermission("users.remove_role"), handler.RemoveRoleFromUser)
		usersGroup.GET("/:userId/permissions-summary", authMiddleware.RequirePermission("users.read_permissions"), handler.GetUserPermissionsSummary)
	}

	// Permission checking endpoint
	authGroup.POST("/check-permission", handler.CheckPermission) // A more general endpoint, might not need admin role
}
