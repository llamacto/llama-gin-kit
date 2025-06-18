package organization

import (
	"context"
	"errors"
)

// Team methods implementation

// CreateTeam adds a new team
func (s *OrganizationServiceImpl) CreateTeam(ctx context.Context, team *Team) error {
	// Verify organization exists before creating team
	_, err := s.GetOrganization(ctx, team.OrganizationID)
	if err != nil {
		return errors.New("organization not found")
	}
	
	return s.repo.CreateTeam(ctx, team)
}

// UpdateTeam updates an existing team
func (s *OrganizationServiceImpl) UpdateTeam(ctx context.Context, team *Team) error {
	// Check if team exists
	existingTeam, err := s.repo.GetTeam(ctx, team.ID)
	if err != nil {
		return errors.New("team not found")
	}
	
	// Prevent change of organization ID
	if team.OrganizationID != existingTeam.OrganizationID {
		return errors.New("cannot change team's organization")
	}
	
	return s.repo.UpdateTeam(ctx, team)
}

// DeleteTeam removes a team by ID
func (s *OrganizationServiceImpl) DeleteTeam(ctx context.Context, id uint) error {
	// Check if team exists
	_, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		return errors.New("team not found")
	}
	
	return s.repo.DeleteTeam(ctx, id)
}

// GetTeam retrieves a team by ID
func (s *OrganizationServiceImpl) GetTeam(ctx context.Context, id uint) (*Team, error) {
	return s.repo.GetTeam(ctx, id)
}

// ListTeams retrieves teams for an organization with pagination
func (s *OrganizationServiceImpl) ListTeams(ctx context.Context, orgID uint, page, pageSize int) ([]*Team, int64, error) {
	// Verify organization exists
	_, err := s.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, 0, errors.New("organization not found")
	}
	
	return s.repo.ListTeams(ctx, orgID, page, pageSize)
}
