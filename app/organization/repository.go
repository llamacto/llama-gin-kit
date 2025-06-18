package organization

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Repository interface for organization data access
type Repository interface {
	CreateOrganization(ctx context.Context, org *Organization) error
	UpdateOrganization(ctx context.Context, org *Organization) error
	DeleteOrganization(ctx context.Context, id uint) error
	GetOrganization(ctx context.Context, id uint) (*Organization, error)
	ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error)
	GetOrganizationsByUserID(ctx context.Context, userID uint) ([]*Organization, error)
}

// repository implementation of Repository
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new organization repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateOrganization adds a new organization
func (r *repository) CreateOrganization(ctx context.Context, org *Organization) error {
	// Debug: Print organization data before saving
	fmt.Printf("Creating organization: %+v\n", org)
	err := r.db.WithContext(ctx).Create(org).Error
	if err != nil {
		fmt.Printf("Error creating organization: %v\n", err)
	}
	return err
}

// UpdateOrganization updates an existing organization
func (r *repository) UpdateOrganization(ctx context.Context, org *Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

// DeleteOrganization removes an organization by ID
func (r *repository) DeleteOrganization(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Organization{}, id).Error
}

// GetOrganization retrieves an organization by ID
func (r *repository) GetOrganization(ctx context.Context, id uint) (*Organization, error) {
	var org Organization
	if err := r.db.WithContext(ctx).First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// ListOrganizations retrieves organizations with pagination
func (r *repository) ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error) {
	var orgs []*Organization
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&Organization{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// GetOrganizationsByUserID retrieves all organizations for a user
func (r *repository) GetOrganizationsByUserID(ctx context.Context, userID uint) ([]*Organization, error) {
	var orgs []*Organization

	err := r.db.WithContext(ctx).
		Joins("JOIN organization_members ON organizations.id = organization_members.organization_id").
		Where("organization_members.user_id = ? AND organization_members.deleted_at IS NULL", userID).
		Find(&orgs).Error

	if err != nil {
		return nil, err
	}
	return orgs, nil
}
