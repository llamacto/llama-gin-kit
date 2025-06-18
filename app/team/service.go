package team

import (
	"fmt"
	"time"
)

// Service defines the interface for team business logic
type Service interface {
	CreateTeam(req *CreateTeamRequest, createdBy uint) (*TeamResponse, error)
	GetTeamByID(id uint) (*TeamResponse, error)
	GetTeamsByOrganization(organizationID uint, page, pageSize int) (*TeamListResponse, error)
	UpdateTeam(id uint, req *UpdateTeamRequest) (*TeamResponse, error)
	DeleteTeam(id uint) error
	GetTeamHierarchy(teamID uint) (*TeamHierarchyResponse, error)
	GetTeamStats(teamID uint) (*TeamWithStats, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new team service instance
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateTeam creates a new team
func (s *service) CreateTeam(req *CreateTeamRequest, createdBy uint) (*TeamResponse, error) {
	// Check if team name already exists in the organization
	exists, err := s.repo.CheckNameExists(req.Name, req.OrganizationID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check team name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("team name '%s' already exists in this organization", req.Name)
	}

	// Create team model
	team := &Team{
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    req.Description,
		OrganizationID: req.OrganizationID,
		ParentTeamID:   req.ParentTeamID,
		// Settings:       req.Settings, // Temporarily disabled
		Status:    1, // Active by default
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	err = s.repo.Create(team)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return s.convertToTeamResponse(team, 0), nil
}

// GetTeamByID retrieves a team by its ID
func (s *service) GetTeamByID(id uint) (*TeamResponse, error) {
	team, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Get team stats
	stats, err := s.repo.GetTeamStats(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team stats: %w", err)
	}

	return s.convertToTeamResponse(team, stats.MemberCount), nil
}

// GetTeamsByOrganization retrieves teams by organization ID with pagination
func (s *service) GetTeamsByOrganization(organizationID uint, page, pageSize int) (*TeamListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	teams, total, err := s.repo.GetByOrganizationID(organizationID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	// Convert to response format
	var teamResponses []TeamResponse
	for _, team := range teams {
		// Get member count for each team
		stats, err := s.repo.GetTeamStats(team.ID)
		memberCount := int64(0)
		if err == nil && stats != nil {
			memberCount = stats.MemberCount
		}

		teamResponses = append(teamResponses, *s.convertToTeamResponse(&team, memberCount))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &TeamListResponse{
		Teams:      teamResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateTeam updates a team
func (s *service) UpdateTeam(id uint, req *UpdateTeamRequest) (*TeamResponse, error) {
	// Check if team exists
	team, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Name != "" {
		// Check if new name already exists (excluding current team)
		exists, err := s.repo.CheckNameExists(req.Name, team.OrganizationID, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check team name existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("team name '%s' already exists in this organization", req.Name)
		}
		updates["name"] = req.Name
	}

	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ParentTeamID != nil {
		updates["parent_team_id"] = req.ParentTeamID
	}
	// if req.Settings != "" {
	//	updates["settings"] = req.Settings
	// } // Temporarily disabled
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	updates["updated_at"] = time.Now()

	// Update team
	err = s.repo.Update(id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	// Return updated team
	return s.GetTeamByID(id)
}

// DeleteTeam deletes a team
func (s *service) DeleteTeam(id uint) error {
	// Check if team exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Check if team has children (prevent deletion if has children)
	children, err := s.repo.GetByParentTeamID(id)
	if err != nil {
		return fmt.Errorf("failed to check team children: %w", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("cannot delete team with child teams")
	}

	// Delete team
	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}

// GetTeamHierarchy retrieves team hierarchy
func (s *service) GetTeamHierarchy(teamID uint) (*TeamHierarchyResponse, error) {
	hierarchy, err := s.repo.GetHierarchy(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team hierarchy: %w", err)
	}

	response := &TeamHierarchyResponse{
		Team: *s.convertToTeamResponse(&hierarchy.Team, 0),
	}

	if hierarchy.Parent != nil {
		parentResponse := s.convertToTeamResponse(hierarchy.Parent, 0)
		response.Parent = parentResponse
	}

	if len(hierarchy.Children) > 0 {
		for _, child := range hierarchy.Children {
			response.Children = append(response.Children, *s.convertToTeamResponse(&child, 0))
		}
	}

	return response, nil
}

// GetTeamStats retrieves team statistics
func (s *service) GetTeamStats(teamID uint) (*TeamWithStats, error) {
	return s.repo.GetTeamStats(teamID)
}

// convertToTeamResponse converts Team model to TeamResponse
func (s *service) convertToTeamResponse(team *Team, memberCount int64) *TeamResponse {
	return &TeamResponse{
		ID:             team.ID,
		Name:           team.Name,
		DisplayName:    team.DisplayName,
		Description:    team.Description,
		OrganizationID: team.OrganizationID,
		ParentTeamID:   team.ParentTeamID,
		// Settings:       team.Settings, // Temporarily disabled
		Status:      team.Status,
		MemberCount: memberCount,
		CreatedAt:   team.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   team.UpdatedAt.Format(time.RFC3339),
	}
}
