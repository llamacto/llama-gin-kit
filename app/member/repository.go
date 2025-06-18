package member

import (
	"gorm.io/gorm"
)

// Repository defines the interface for member data operations
type Repository interface {
	Create(member *Member) error
	GetByID(id uint) (*Member, error)
	GetByUserAndOrganization(userID, organizationID uint) (*Member, error)
	GetByOrganizationID(organizationID uint, page, pageSize int) ([]MemberWithDetails, int64, error)
	GetByTeamID(teamID uint, page, pageSize int) ([]MemberWithDetails, int64, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	GetMemberStats(organizationID uint) (*MemberStatsResponse, error)
	CheckMemberExists(userID, organizationID uint) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new member repository instance
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new member
func (r *repository) Create(member *Member) error {
	return r.db.Create(member).Error
}

// GetByID retrieves a member by its ID
func (r *repository) GetByID(id uint) (*Member, error) {
	var member Member
	err := r.db.First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetByUserAndOrganization retrieves a member by user ID and organization ID
func (r *repository) GetByUserAndOrganization(userID, organizationID uint) (*Member, error) {
	var member Member
	err := r.db.Where("user_id = ? AND organization_id = ?", userID, organizationID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetByOrganizationID retrieves members by organization ID with pagination and detailed info
func (r *repository) GetByOrganizationID(organizationID uint, page, pageSize int) ([]MemberWithDetails, int64, error) {
	var members []MemberWithDetails
	var total int64

	// Count total records
	err := r.db.Table("organization_members").
		Where("organization_id = ? AND deleted_at IS NULL", organizationID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results with joins
	offset := (page - 1) * pageSize
	err = r.db.Table("organization_members as om").
		Select(`
			om.id, om.user_id, om.organization_id, om.team_id, om.role_id,
			om.status, om.joined_at, om.invited_by, om.created_at, om.updated_at,
			u.name as user_name, u.email as user_email, u.nickname as user_nickname, u.avatar as user_avatar,
			o.name as organization_name,
			t.name as team_name,
			r.name as role_name, r.display_name as role_display_name
		`).
		Joins("LEFT JOIN users u ON om.user_id = u.id").
		Joins("LEFT JOIN organizations o ON om.organization_id = o.id").
		Joins("LEFT JOIN teams t ON om.team_id = t.id").
		Joins("LEFT JOIN organization_roles r ON om.role_id = r.id").
		Where("om.organization_id = ? AND om.deleted_at IS NULL", organizationID).
		Offset(offset).
		Limit(pageSize).
		Scan(&members).Error

	return members, total, err
}

// GetByTeamID retrieves members by team ID with pagination and detailed info
func (r *repository) GetByTeamID(teamID uint, page, pageSize int) ([]MemberWithDetails, int64, error) {
	var members []MemberWithDetails
	var total int64

	// Count total records
	err := r.db.Table("organization_members").
		Where("team_id = ? AND deleted_at IS NULL", teamID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results with joins
	offset := (page - 1) * pageSize
	err = r.db.Table("organization_members as om").
		Select(`
			om.id, om.user_id, om.organization_id, om.team_id, om.role_id,
			om.status, om.joined_at, om.invited_by, om.created_at, om.updated_at,
			u.name as user_name, u.email as user_email, u.nickname as user_nickname, u.avatar as user_avatar,
			o.name as organization_name,
			t.name as team_name,
			r.name as role_name, r.display_name as role_display_name
		`).
		Joins("LEFT JOIN users u ON om.user_id = u.id").
		Joins("LEFT JOIN organizations o ON om.organization_id = o.id").
		Joins("LEFT JOIN teams t ON om.team_id = t.id").
		Joins("LEFT JOIN organization_roles r ON om.role_id = r.id").
		Where("om.team_id = ? AND om.deleted_at IS NULL", teamID).
		Offset(offset).
		Limit(pageSize).
		Scan(&members).Error

	return members, total, err
}

// Update updates a member by ID
func (r *repository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&Member{}).Where("id = ?", id).Updates(updates).Error
}

// Delete soft deletes a member by ID
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Member{}, id).Error
}

// GetMemberStats retrieves member statistics for an organization
func (r *repository) GetMemberStats(organizationID uint) (*MemberStatsResponse, error) {
	stats := &MemberStatsResponse{}

	// Total members
	err := r.db.Table("organization_members").
		Where("organization_id = ? AND deleted_at IS NULL", organizationID).
		Count(&stats.TotalMembers).Error
	if err != nil {
		return nil, err
	}

	// Active members
	err = r.db.Table("organization_members").
		Where("organization_id = ? AND status = 1 AND deleted_at IS NULL", organizationID).
		Count(&stats.ActiveMembers).Error
	if err != nil {
		return nil, err
	}

	// Pending invites
	err = r.db.Table("organization_invitations").
		Where("organization_id = ? AND status = 0 AND deleted_at IS NULL", organizationID).
		Count(&stats.PendingInvites).Error
	if err != nil {
		return nil, err
	}

	// Disabled members
	err = r.db.Table("organization_members").
		Where("organization_id = ? AND status = 2 AND deleted_at IS NULL", organizationID).
		Count(&stats.DisabledMembers).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// CheckMemberExists checks if a user is already a member of the organization
func (r *repository) CheckMemberExists(userID, organizationID uint) (bool, error) {
	var count int64
	err := r.db.Table("organization_members").
		Where("user_id = ? AND organization_id = ? AND deleted_at IS NULL", userID, organizationID).
		Count(&count).Error
	return count > 0, err
}
