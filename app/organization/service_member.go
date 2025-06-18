package organization

import (
	"context"
	"errors"
	"time"
)

// Member methods implementation

// AddMember adds a new member to an organization
func (s *OrganizationServiceImpl) AddMember(ctx context.Context, member *Member) error {
	// Verify organization exists
	_, err := s.GetOrganization(ctx, member.OrganizationID)
	if err != nil {
		return errors.New("organization not found")
	}
	
	// Verify role exists
	_, err = s.GetRole(ctx, member.RoleID)
	if err != nil {
		return errors.New("role not found")
	}
	
	// Check if member already exists
	existingMember, err := s.repo.GetMemberByUserAndOrg(ctx, member.UserID, member.OrganizationID)
	if err == nil && existingMember != nil {
		return errors.New("user is already a member of this organization")
	}
	
	// Set joined time
	member.JoinedAt = time.Now()
	
	return s.repo.AddMember(ctx, member)
}

// UpdateMember updates an existing member
func (s *OrganizationServiceImpl) UpdateMember(ctx context.Context, member *Member) error {
	// Check if member exists
	existingMember, err := s.GetMember(ctx, member.ID)
	if err != nil {
		return errors.New("member not found")
	}
	
	// Prevent change of organization or user
	if member.OrganizationID != existingMember.OrganizationID || 
	   member.UserID != existingMember.UserID {
		return errors.New("cannot change organization or user ID of a member")
	}
	
	// Verify role exists
	_, err = s.GetRole(ctx, member.RoleID)
	if err != nil {
		return errors.New("role not found")
	}
	
	return s.repo.UpdateMember(ctx, member)
}

// RemoveMember removes a member by ID
func (s *OrganizationServiceImpl) RemoveMember(ctx context.Context, id uint) error {
	// Check if member exists
	_, err := s.GetMember(ctx, id)
	if err != nil {
		return errors.New("member not found")
	}
	
	return s.repo.RemoveMember(ctx, id)
}

// GetMember retrieves a member by ID
func (s *OrganizationServiceImpl) GetMember(ctx context.Context, id uint) (*Member, error) {
	return s.repo.GetMember(ctx, id)
}

// ListMembers retrieves members for an organization with pagination
func (s *OrganizationServiceImpl) ListMembers(ctx context.Context, orgID uint, page, pageSize int) ([]*Member, int64, error) {
	// Verify organization exists
	_, err := s.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, 0, errors.New("organization not found")
	}
	
	return s.repo.ListMembers(ctx, orgID, page, pageSize)
}

// ListTeamMembers retrieves members for a team with pagination
func (s *OrganizationServiceImpl) ListTeamMembers(ctx context.Context, teamID uint, page, pageSize int) ([]*Member, int64, error) {
	// Verify team exists
	_, err := s.GetTeam(ctx, teamID)
	if err != nil {
		return nil, 0, errors.New("team not found")
	}
	
	return s.repo.ListTeamMembers(ctx, teamID, page, pageSize)
}
