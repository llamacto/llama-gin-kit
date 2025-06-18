package team

import (
	"gorm.io/gorm"
)

// Repository defines the interface for team data operations
type Repository interface {
	Create(team *Team) error
	GetByID(id uint) (*Team, error)
	GetByOrganizationID(organizationID uint, page, pageSize int) ([]Team, int64, error)
	GetByParentTeamID(parentTeamID uint) ([]Team, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	GetHierarchy(teamID uint) (*TeamHierarchy, error)
	GetTeamStats(teamID uint) (*TeamWithStats, error)
	CheckNameExists(name string, organizationID uint, excludeID *uint) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new team repository instance
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new team
func (r *repository) Create(team *Team) error {
	return r.db.Create(team).Error
}

// GetByID retrieves a team by its ID
func (r *repository) GetByID(id uint) (*Team, error) {
	var team Team
	err := r.db.First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// GetByOrganizationID retrieves teams by organization ID with pagination
func (r *repository) GetByOrganizationID(organizationID uint, page, pageSize int) ([]Team, int64, error) {
	var teams []Team
	var total int64

	query := r.db.Where("organization_id = ?", organizationID)

	// Count total records
	err := query.Model(&Team{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&teams).Error
	if err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

// GetByParentTeamID retrieves teams by parent team ID
func (r *repository) GetByParentTeamID(parentTeamID uint) ([]Team, error) {
	var teams []Team
	err := r.db.Where("parent_team_id = ?", parentTeamID).Find(&teams).Error
	return teams, err
}

// Update updates a team by ID
func (r *repository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&Team{}).Where("id = ?", id).Updates(updates).Error
}

// Delete soft deletes a team by ID
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Team{}, id).Error
}

// GetHierarchy retrieves team hierarchy (parent and children)
func (r *repository) GetHierarchy(teamID uint) (*TeamHierarchy, error) {
	var team Team
	err := r.db.First(&team, teamID).Error
	if err != nil {
		return nil, err
	}

	hierarchy := &TeamHierarchy{
		Team: team,
	}

	// Get parent team if exists
	if team.ParentTeamID != nil {
		var parent Team
		if err := r.db.First(&parent, *team.ParentTeamID).Error; err == nil {
			hierarchy.Parent = &parent
		}
	}

	// Get children teams
	var children []Team
	if err := r.db.Where("parent_team_id = ?", teamID).Find(&children).Error; err == nil {
		hierarchy.Children = children
	}

	return hierarchy, nil
}

// GetTeamStats retrieves team with member count statistics
func (r *repository) GetTeamStats(teamID uint) (*TeamWithStats, error) {
	var team Team
	err := r.db.First(&team, teamID).Error
	if err != nil {
		return nil, err
	}

	var memberCount int64
	err = r.db.Table("organization_members").
		Where("team_id = ? AND deleted_at IS NULL", teamID).
		Count(&memberCount).Error
	if err != nil {
		return nil, err
	}

	return &TeamWithStats{
		Team:        team,
		MemberCount: memberCount,
	}, nil
}

// CheckNameExists checks if a team name already exists in the organization
func (r *repository) CheckNameExists(name string, organizationID uint, excludeID *uint) (bool, error) {
	query := r.db.Where("name = ? AND organization_id = ?", name, organizationID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Model(&Team{}).Count(&count).Error
	return count > 0, err
}
