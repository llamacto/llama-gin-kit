package organization

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/llamacto/llama-gin-kit/app/user"
	"gorm.io/gorm"
)

// OrganizationService interface for organization business logic
type OrganizationService interface {
	// Organization methods
	CreateOrganization(ctx context.Context, org *Organization, userID uint) error
	UpdateOrganization(ctx context.Context, org *Organization) error
	DeleteOrganization(ctx context.Context, id uint) error
	GetOrganization(ctx context.Context, id uint) (*Organization, error)
	ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error)
	GetUserOrganizations(ctx context.Context, userID uint) ([]*Organization, error)
	
	// Team methods
	CreateTeam(ctx context.Context, team *Team) error
	UpdateTeam(ctx context.Context, team *Team) error
	DeleteTeam(ctx context.Context, id uint) error
	GetTeam(ctx context.Context, id uint) (*Team, error)
	ListTeams(ctx context.Context, orgID uint, page, pageSize int) ([]*Team, int64, error)
	
	// Member methods
	AddMember(ctx context.Context, member *Member) error
	UpdateMember(ctx context.Context, member *Member) error
	RemoveMember(ctx context.Context, id uint) error
	GetMember(ctx context.Context, id uint) (*Member, error)
	ListMembers(ctx context.Context, orgID uint, page, pageSize int) ([]*Member, int64, error)
	ListTeamMembers(ctx context.Context, teamID uint, page, pageSize int) ([]*Member, int64, error)
	
	// Role methods
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRole(ctx context.Context, id uint) (*Role, error)
	ListRoles(ctx context.Context, orgID uint, page, pageSize int) ([]*Role, int64, error)
	
	// Permission methods
	CheckPermission(ctx context.Context, userID uint, orgID uint, permission string) (bool, error)
	
	// Invitation methods
	InviteMember(ctx context.Context, invitation *Invitation) error
	ProcessInvitation(ctx context.Context, token string, userID uint) error
	CancelInvitation(ctx context.Context, id uint) error
	GetInvitation(ctx context.Context, id uint) (*Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*Invitation, error)
	ListInvitations(ctx context.Context, orgID uint, page, pageSize int) ([]*Invitation, int64, error)
}

// OrganizationServiceImpl implementation of OrganizationService
type OrganizationServiceImpl struct {
	repo OrganizationRepository
	userService user.UserService
	db *gorm.DB
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(repo OrganizationRepository, userService user.UserService, db *gorm.DB) OrganizationService {
	return &OrganizationServiceImpl{
		repo: repo,
		userService: userService,
		db: db,
	}
}

// GenerateToken creates a secure random token for invitations
func GenerateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Organization methods implementation

// CreateOrganization adds a new organization and adds the creator as admin
func (s *OrganizationServiceImpl) CreateOrganization(ctx context.Context, org *Organization, userID uint) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		orgRepo := NewOrganizationRepository(tx)
		
		// Create the organization
		if err := orgRepo.CreateOrganization(ctx, org); err != nil {
			return err
		}
		
		// Get or create admin role
		var adminRole Role
		if err := tx.Where("name = ? AND (organization_id = ? OR organization_id IS NULL)", "admin", org.ID).
			First(&adminRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create admin role if not exists
				adminRole = Role{
					Name:           "admin",
					DisplayName:    "Administrator",
					Description:    "Full administrative access",
					OrganizationID: &org.ID,
					Permissions:    `{"*": true}`, // All permissions
					IsDefault:      false,
				}
				if err := orgRepo.CreateRole(ctx, &adminRole); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		
		// Add creator as admin member
		member := &Member{
			UserID:         userID,
			OrganizationID: org.ID,
			RoleID:         adminRole.ID,
			Status:         1, // Active
			JoinedAt:       time.Now(),
			InvitedBy:      userID,
		}
		
		if err := orgRepo.AddMember(ctx, member); err != nil {
			return err
		}
		
		return nil
	})
	
	return err
}

// UpdateOrganization updates an existing organization
func (s *OrganizationServiceImpl) UpdateOrganization(ctx context.Context, org *Organization) error {
	return s.repo.UpdateOrganization(ctx, org)
}

// DeleteOrganization removes an organization by ID
func (s *OrganizationServiceImpl) DeleteOrganization(ctx context.Context, id uint) error {
	return s.repo.DeleteOrganization(ctx, id)
}

// GetOrganization retrieves an organization by ID
func (s *OrganizationServiceImpl) GetOrganization(ctx context.Context, id uint) (*Organization, error) {
	return s.repo.GetOrganization(ctx, id)
}

// ListOrganizations retrieves organizations with pagination
func (s *OrganizationServiceImpl) ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error) {
	return s.repo.ListOrganizations(ctx, page, pageSize)
}

// GetUserOrganizations retrieves all organizations for a user
func (s *OrganizationServiceImpl) GetUserOrganizations(ctx context.Context, userID uint) ([]*Organization, error) {
	return s.repo.GetOrganizationsByUserID(ctx, userID)
}
