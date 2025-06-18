package organization

import (
	"context"

	"github.com/llamacto/llama-gin-kit/app/user"
	"gorm.io/gorm"
)

// Service interface for organization business logic
type Service interface {
	CreateOrganization(ctx context.Context, org *Organization, userID uint) error
	UpdateOrganization(ctx context.Context, org *Organization) error
	DeleteOrganization(ctx context.Context, id uint) error
	GetOrganization(ctx context.Context, id uint) (*Organization, error)
	ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error)
	GetUserOrganizations(ctx context.Context, userID uint) ([]*Organization, error)
	GetOrganizationStats(ctx context.Context, id uint) (*OrganizationStats, error)
}

// service implementation of Service
type service struct {
	repo        Repository
	userService user.UserService
	db          *gorm.DB
}

// NewService creates a new organization service
func NewService(repo Repository, userService user.UserService, db *gorm.DB) Service {
	return &service{
		repo:        repo,
		userService: userService,
		db:          db,
	}
}

// CreateOrganization adds a new organization
func (s *service) CreateOrganization(ctx context.Context, org *Organization, userID uint) error {
	return s.repo.CreateOrganization(ctx, org)
}

// UpdateOrganization updates an existing organization
func (s *service) UpdateOrganization(ctx context.Context, org *Organization) error {
	return s.repo.UpdateOrganization(ctx, org)
}

// DeleteOrganization removes an organization by ID
func (s *service) DeleteOrganization(ctx context.Context, id uint) error {
	return s.repo.DeleteOrganization(ctx, id)
}

// GetOrganization retrieves an organization by ID
func (s *service) GetOrganization(ctx context.Context, id uint) (*Organization, error) {
	return s.repo.GetOrganization(ctx, id)
}

// ListOrganizations retrieves organizations with pagination
func (s *service) ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error) {
	return s.repo.ListOrganizations(ctx, page, pageSize)
}

// GetUserOrganizations retrieves all organizations for a user
func (s *service) GetUserOrganizations(ctx context.Context, userID uint) ([]*Organization, error) {
	return s.repo.GetOrganizationsByUserID(ctx, userID)
}

// GetOrganizationStats retrieves organization statistics
func (s *service) GetOrganizationStats(ctx context.Context, id uint) (*OrganizationStats, error) {
	org, err := s.repo.GetOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	stats := &OrganizationStats{
		Organization: *org,
	}

	// Get member count
	err = s.db.Table("organization_members").
		Where("organization_id = ? AND deleted_at IS NULL", id).
		Count(&stats.MemberCount).Error
	if err != nil {
		return nil, err
	}

	// Get team count
	err = s.db.Table("teams").
		Where("organization_id = ? AND deleted_at IS NULL", id).
		Count(&stats.TeamCount).Error
	if err != nil {
		return nil, err
	}

	// Get role count
	err = s.db.Table("organization_roles").
		Where("organization_id = ? AND deleted_at IS NULL", id).
		Count(&stats.RoleCount).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}
