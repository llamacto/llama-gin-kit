package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/apikey"
	"github.com/llamacto/llama-gin-kit/app/organization"
	apikeyMiddleware "github.com/llamacto/llama-gin-kit/middleware"
)

// RegisterOrganizationRoutes registers organization routes
func RegisterOrganizationRoutes(router *gin.RouterGroup, handler *organization.Handler, apiKeyService apikey.Service) {
	// Routes that require authentication
	authRouter := router.Group("")
	authRouter.Use(apikeyMiddleware.CombinedAuth(apiKeyService))

	// Organization endpoints
	orgRouter := authRouter.Group("/organizations")
	orgRouter.POST("", handler.CreateOrganization)
	orgRouter.GET("", handler.ListOrganizations)
	orgRouter.GET("/me", handler.GetMyOrganizations)
	orgRouter.GET("/:id", handler.GetOrganization)
	orgRouter.PUT("/:id", handler.UpdateOrganization)
	orgRouter.DELETE("/:id", handler.DeleteOrganization)
	orgRouter.GET("/:id/teams", handler.ListTeams)
	orgRouter.GET("/:id/members", handler.ListMembers)
	orgRouter.GET("/:id/roles", handler.ListRoles)
	orgRouter.GET("/:id/invitations", handler.ListInvitations)
	
	// Team endpoints
	teamRouter := authRouter.Group("/teams")
	teamRouter.POST("", handler.CreateTeam)
	teamRouter.GET("/:id", handler.GetTeam)
	teamRouter.PUT("/:id", handler.UpdateTeam)
	teamRouter.DELETE("/:id", handler.DeleteTeam)
	
	// Member endpoints
	memberRouter := authRouter.Group("/members")
	memberRouter.POST("", handler.AddMember)
	memberRouter.GET("/:id", handler.GetMember)
	memberRouter.PUT("/:id", handler.UpdateMember)
	memberRouter.DELETE("/:id", handler.RemoveMember)
	
	// Role endpoints
	roleRouter := authRouter.Group("/roles")
	roleRouter.POST("", handler.CreateRole)
	roleRouter.GET("/:id", handler.GetRole)
	roleRouter.PUT("/:id", handler.UpdateRole)
	roleRouter.DELETE("/:id", handler.DeleteRole)
	
	// Invitation endpoints
	invitationRouter := authRouter.Group("/invitations")
	invitationRouter.POST("", handler.CreateInvitation)
	invitationRouter.GET("/:id", handler.GetInvitation)
	invitationRouter.DELETE("/:id", handler.CancelInvitation)
	invitationRouter.POST("/accept", handler.AcceptInvitation)
	invitationRouter.GET("/token/:token", handler.GetInvitationByToken)
	
	// Permission endpoints
	permissionRouter := authRouter.Group("/permissions")
	permissionRouter.POST("/check", handler.CheckPermission)
}
