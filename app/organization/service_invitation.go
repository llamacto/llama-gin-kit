package organization

import (
	"context"
	"errors"
	"time"
)

// Invitation methods implementation

// InviteMember sends an invitation to join an organization
func (s *OrganizationServiceImpl) InviteMember(ctx context.Context, invitation *Invitation) error {
	// Verify organization exists
	_, err := s.GetOrganization(ctx, invitation.OrganizationID)
	if err != nil {
		return errors.New("organization not found")
	}
	
	// Verify role exists
	_, err = s.GetRole(ctx, invitation.RoleID)
	if err != nil {
		return errors.New("role not found")
	}
	
	// If team is specified, verify it exists
	if invitation.TeamID != nil && *invitation.TeamID > 0 {
		_, err = s.GetTeam(ctx, *invitation.TeamID)
		if err != nil {
			return errors.New("team not found")
		}
	}
	
	// Generate unique invitation token
	token, err := GenerateToken(32)
	if err != nil {
		return errors.New("failed to generate invitation token")
	}
	invitation.Token = token
	
	// Set expiration time (default to 7 days)
	invitation.ExpiresAt = time.Now().AddDate(0, 0, 7)
	
	// Set initial status to pending
	invitation.Status = 0
	
	return s.repo.CreateInvitation(ctx, invitation)
}

// ProcessInvitation accepts or rejects an invitation
func (s *OrganizationServiceImpl) ProcessInvitation(ctx context.Context, token string, userID uint) error {
	// Get the invitation by token
	invitation, err := s.repo.GetInvitationByToken(ctx, token)
	if err != nil {
		return errors.New("invitation not found")
	}
	
	// Check if invitation is expired
	if invitation.ExpiresAt.Before(time.Now()) {
		invitation.Status = 3 // Expired
		_ = s.repo.UpdateInvitation(ctx, invitation)
		return errors.New("invitation has expired")
	}
	
	// Check if invitation is already processed
	if invitation.Status != 0 {
		return errors.New("invitation has already been processed")
	}
	
	// Update invitation status to accepted
	invitation.Status = 1 // Accepted
	if err := s.repo.UpdateInvitation(ctx, invitation); err != nil {
		return err
	}
	
	// Create a new member entry
	member := &Member{
		UserID:         userID,
		OrganizationID: invitation.OrganizationID,
		TeamID:         invitation.TeamID,
		RoleID:         invitation.RoleID,
		Status:         1, // Active
		JoinedAt:       time.Now(),
		InvitedBy:      invitation.InvitedBy,
	}
	
	return s.repo.AddMember(ctx, member)
}

// CancelInvitation cancels a pending invitation
func (s *OrganizationServiceImpl) CancelInvitation(ctx context.Context, id uint) error {
	// Get the invitation
	invitation, err := s.repo.GetInvitation(ctx, id)
	if err != nil {
		return errors.New("invitation not found")
	}
	
	// Check if invitation can be cancelled
	if invitation.Status != 0 {
		return errors.New("only pending invitations can be cancelled")
	}
	
	// Update invitation status to rejected
	invitation.Status = 2 // Rejected
	return s.repo.UpdateInvitation(ctx, invitation)
}

// GetInvitation retrieves an invitation by ID
func (s *OrganizationServiceImpl) GetInvitation(ctx context.Context, id uint) (*Invitation, error) {
	return s.repo.GetInvitation(ctx, id)
}

// GetInvitationByToken retrieves an invitation by token
func (s *OrganizationServiceImpl) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	return s.repo.GetInvitationByToken(ctx, token)
}

// ListInvitations retrieves invitations for an organization with pagination
func (s *OrganizationServiceImpl) ListInvitations(ctx context.Context, orgID uint, page, pageSize int) ([]*Invitation, int64, error) {
	// Verify organization exists
	_, err := s.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, 0, errors.New("organization not found")
	}
	
	return s.repo.ListInvitations(ctx, orgID, page, pageSize)
}
