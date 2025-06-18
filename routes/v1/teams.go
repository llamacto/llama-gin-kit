package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/app/team"
	"github.com/llamacto/llama-gin-kit/pkg/database"
	pkgmiddleware "github.com/llamacto/llama-gin-kit/pkg/middleware"
)

// TeamRoutes sets up team-related routes
func TeamRoutes(router *gin.RouterGroup) {
	// Initialize team dependencies
	teamRepo := team.NewRepository(database.DB)
	teamService := team.NewService(teamRepo)
	teamHandler := team.NewHandler(teamService)

	// Team routes group
	teams := router.Group("/teams")
	teams.Use(pkgmiddleware.JWTAuth()) // Require authentication for all team operations
	{
		teams.POST("", teamHandler.CreateTeam)                    // Create team
		teams.GET("/:id", teamHandler.GetTeam)                    // Get team by ID
		teams.PUT("/:id", teamHandler.UpdateTeam)                 // Update team
		teams.DELETE("/:id", teamHandler.DeleteTeam)              // Delete team
		teams.GET("/:id/hierarchy", teamHandler.GetTeamHierarchy) // Get team hierarchy
	}

	// Organization-specific team routes - moved to avoid route conflicts
	orgTeams := router.Group("/org-teams")
	orgTeams.Use(pkgmiddleware.JWTAuth())
	{
		orgTeams.GET("/:organization_id", teamHandler.GetTeamsByOrganization) // Get organization teams
	}
}
