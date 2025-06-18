package organization

import (
	"context"

	"gorm.io/gorm"
)

// OrganizationRepository interface for organization data access
type OrganizationRepository interface {
	// Organization methods
	CreateOrganization(ctx context.Context, org *Organization) error
	UpdateOrganization(ctx context.Context, org *Organization) error
	DeleteOrganization(ctx context.Context, id uint) error
	GetOrganization(ctx context.Context, id uint) (*Organization, error)
	ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error)
	GetOrganizationsByUserID(ctx context.Context, userID uint) ([]*Organization, error)
	
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
	GetMemberByUserAndOrg(ctx context.Context, userID, orgID uint) (*Member, error)
	ListMembers(ctx context.Context, orgID uint, page, pageSize int) ([]*Member, int64, error)
	ListTeamMembers(ctx context.Context, teamID uint, page, pageSize int) ([]*Member, int64, error)
	
	// Role methods
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRole(ctx context.Context, id uint) (*Role, error)
	ListRoles(ctx context.Context, orgID uint, page, pageSize int) ([]*Role, int64, error)
	
	// Invitation methods
	CreateInvitation(ctx context.Context, invitation *Invitation) error
	UpdateInvitation(ctx context.Context, invitation *Invitation) error
	DeleteInvitation(ctx context.Context, id uint) error
	GetInvitation(ctx context.Context, id uint) (*Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*Invitation, error)
	ListInvitations(ctx context.Context, orgID uint, page, pageSize int) ([]*Invitation, int64, error)
}

// OrganizationRepositoryImpl implementation of OrganizationRepository
type OrganizationRepositoryImpl struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &OrganizationRepositoryImpl{db: db}
}

// Organization methods implementation

// CreateOrganization adds a new organization
func (r *OrganizationRepositoryImpl) CreateOrganization(ctx context.Context, org *Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

// UpdateOrganization updates an existing organization
func (r *OrganizationRepositoryImpl) UpdateOrganization(ctx context.Context, org *Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

// DeleteOrganization removes an organization by ID
func (r *OrganizationRepositoryImpl) DeleteOrganization(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Organization{}, id).Error
}

// GetOrganization retrieves an organization by ID
func (r *OrganizationRepositoryImpl) GetOrganization(ctx context.Context, id uint) (*Organization, error) {
	var org Organization
	if err := r.db.WithContext(ctx).First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// ListOrganizations retrieves organizations with pagination
func (r *OrganizationRepositoryImpl) ListOrganizations(ctx context.Context, page, pageSize int) ([]*Organization, int64, error) {
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
func (r *OrganizationRepositoryImpl) GetOrganizationsByUserID(ctx context.Context, userID uint) ([]*Organization, error) {
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

// Team methods implementation

// CreateTeam adds a new team
func (r *OrganizationRepositoryImpl) CreateTeam(ctx context.Context, team *Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

// UpdateTeam updates an existing team
func (r *OrganizationRepositoryImpl) UpdateTeam(ctx context.Context, team *Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

// DeleteTeam removes a team by ID
func (r *OrganizationRepositoryImpl) DeleteTeam(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Team{}, id).Error
}

// GetTeam retrieves a team by ID
func (r *OrganizationRepositoryImpl) GetTeam(ctx context.Context, id uint) (*Team, error) {
	var team Team
	if err := r.db.WithContext(ctx).First(&team, id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

// ListTeams retrieves teams for an organization with pagination
func (r *OrganizationRepositoryImpl) ListTeams(ctx context.Context, orgID uint, page, pageSize int) ([]*Team, int64, error) {
	var teams []*Team
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&Team{}).Where("organization_id = ?", orgID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Offset(offset).Limit(pageSize).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

// Member methods implementation

// AddMember adds a new member to an organization
func (r *OrganizationRepositoryImpl) AddMember(ctx context.Context, member *Member) error {
	return r.db.WithContext(ctx).Create(member).Error
}

// UpdateMember updates an existing member
func (r *OrganizationRepositoryImpl) UpdateMember(ctx context.Context, member *Member) error {
	return r.db.WithContext(ctx).Save(member).Error
}

// RemoveMember removes a member by ID
func (r *OrganizationRepositoryImpl) RemoveMember(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Member{}, id).Error
}

// GetMember retrieves a member by ID
func (r *OrganizationRepositoryImpl) GetMember(ctx context.Context, id uint) (*Member, error) {
	var member Member
	if err := r.db.WithContext(ctx).First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMemberByUserAndOrg retrieves a member by user ID and organization ID
func (r *OrganizationRepositoryImpl) GetMemberByUserAndOrg(ctx context.Context, userID, orgID uint) (*Member, error) {
	var member Member
	if err := r.db.WithContext(ctx).Where("user_id = ? AND organization_id = ?", userID, orgID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// ListMembers retrieves members for an organization with pagination
func (r *OrganizationRepositoryImpl) ListMembers(ctx context.Context, orgID uint, page, pageSize int) ([]*Member, int64, error) {
	var members []*Member
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&Member{}).Where("organization_id = ?", orgID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Offset(offset).Limit(pageSize).Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// ListTeamMembers retrieves members for a team with pagination
func (r *OrganizationRepositoryImpl) ListTeamMembers(ctx context.Context, teamID uint, page, pageSize int) ([]*Member, int64, error) {
	var members []*Member
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&Member{}).Where("team_id = ?", teamID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Offset(offset).Limit(pageSize).Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// Role methods implementation

// CreateRole adds a new role
func (r *OrganizationRepositoryImpl) CreateRole(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// UpdateRole updates an existing role
func (r *OrganizationRepositoryImpl) UpdateRole(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// DeleteRole removes a role by ID
func (r *OrganizationRepositoryImpl) DeleteRole(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Role{}, id).Error
}

// GetRole retrieves a role by ID
func (r *OrganizationRepositoryImpl) GetRole(ctx context.Context, id uint) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles retrieves roles for an organization with pagination
func (r *OrganizationRepositoryImpl) ListRoles(ctx context.Context, orgID uint, page, pageSize int) ([]*Role, int64, error) {
	var roles []*Role
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&Role{})
	if orgID > 0 {
		query = query.Where("organization_id = ? OR organization_id IS NULL", orgID)
	} else {
		query = query.Where("organization_id IS NULL") // System roles only
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Invitation methods implementation

// CreateInvitation adds a new invitation
func (r *OrganizationRepositoryImpl) CreateInvitation(ctx context.Context, invitation *Invitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

// UpdateInvitation updates an existing invitation
func (r *OrganizationRepositoryImpl) UpdateInvitation(ctx context.Context, invitation *Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

// DeleteInvitation removes an invitation by ID
func (r *OrganizationRepositoryImpl) DeleteInvitation(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Invitation{}, id).Error
}

// GetInvitation retrieves an invitation by ID
func (r *OrganizationRepositoryImpl) GetInvitation(ctx context.Context, id uint) (*Invitation, error) {
	var invitation Invitation
	if err := r.db.WithContext(ctx).First(&invitation, id).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetInvitationByToken retrieves an invitation by token
func (r *OrganizationRepositoryImpl) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	var invitation Invitation
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

// ListInvitations retrieves invitations for an organization with pagination
func (r *OrganizationRepositoryImpl) ListInvitations(ctx context.Context, orgID uint, page, pageSize int) ([]*Invitation, int64, error) {
	var invitations []*Invitation
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&Invitation{}).Where("organization_id = ?", orgID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Offset(offset).Limit(pageSize).Find(&invitations).Error; err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}
