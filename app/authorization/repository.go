package authorization

import (
	"fmt"

	"gorm.io/gorm"
)

// Repository interface for authorization data access
type Repository interface {
	// Role operations
	CreateRole(role *Role) error
	GetRoleByID(id uint) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	UpdateRole(role *Role) error
	DeleteRole(id uint) error
	ListRoles(query ListRolesQuery) ([]Role, int64, error)
	GetRoleWithPermissions(id uint) (*Role, error)

	// Permission operations
	CreatePermission(permission *Permission) error
	GetPermissionByID(id uint) (*Permission, error)
	GetPermissionByName(name string) (*Permission, error)
	UpdatePermission(permission *Permission) error
	DeletePermission(id uint) error
	ListPermissions(query ListPermissionsQuery) ([]Permission, int64, error)
	GetPermissionsByIDs(ids []uint) ([]Permission, error)

	// Role-Permission operations
	AssignPermissionsToRole(roleID uint, permissionIDs []uint, grantedBy uint) error
	RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]Permission, error)

	// User-Role operations
	AssignRoleToUser(userRole *UserRole) error
	RemoveRoleFromUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]UserRole, error)
	GetUsersWithRole(roleID uint) ([]UserRole, error)
	CheckUserRole(userID, roleID uint) (bool, error)

	// Organization-Role operations
	AssignOrganizationRole(orgRole *OrganizationRole) error
	RemoveOrganizationRole(userID, organizationID, roleID uint) error
	GetUserOrganizationRoles(userID, organizationID uint) ([]OrganizationRole, error)
	GetOrganizationUsers(organizationID, roleID uint) ([]OrganizationRole, error)

	// Team-Role operations
	AssignTeamRole(teamRole *TeamRole) error
	RemoveTeamRole(userID, teamID, roleID uint) error
	GetUserTeamRoles(userID, teamID uint) ([]TeamRole, error)
	GetTeamUsers(teamID, roleID uint) ([]TeamRole, error)

	// Policy operations
	CreatePolicy(policy *Policy) error
	GetPolicyByID(id uint) (*Policy, error)
	UpdatePolicy(policy *Policy) error
	DeletePolicy(id uint) error
	ListPolicies(query ListQuery) ([]Policy, int64, error)

	// Permission checking operations
	GetUserAllPermissions(userID uint) ([]string, error)
	GetUserOrganizationPermissions(userID, organizationID uint) ([]string, error)
	GetUserTeamPermissions(userID, teamID uint) ([]string, error)
	CheckUserPermission(userID uint, permission string) (bool, error)
	CheckUserOrganizationPermission(userID, organizationID uint, permission string) (bool, error)
	CheckUserTeamPermission(userID, teamID uint, permission string) (bool, error)
}

// repositoryImpl implements the Repository interface
type repositoryImpl struct {
	db *gorm.DB
}

// NewRepository creates a new authorization repository
func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

// ===== Role operations =====

func (r *repositoryImpl) CreateRole(role *Role) error {
	return r.db.Create(role).Error
}

func (r *repositoryImpl) GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repositoryImpl) GetRoleByName(name string) (*Role, error) {
	var role Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repositoryImpl) UpdateRole(role *Role) error {
	return r.db.Save(role).Error
}

func (r *repositoryImpl) DeleteRole(id uint) error {
	return r.db.Delete(&Role{}, id).Error
}

func (r *repositoryImpl) ListRoles(query ListRolesQuery) ([]Role, int64, error) {
	var roles []Role
	var total int64

	db := r.db.Model(&Role{})

	// Apply filters
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR display_name ILIKE ?", searchPattern, searchPattern)
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	if query.Level != nil {
		db = db.Where("level = ?", *query.Level)
	}

	if query.IsSystem != nil {
		db = db.Where("is_system = ?", *query.IsSystem)
	}

	// Count total
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (query.Page - 1) * query.PageSize
	orderBy := fmt.Sprintf("%s %s", query.OrderBy, query.Order)

	err = db.Order(orderBy).Offset(offset).Limit(query.PageSize).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *repositoryImpl) GetRoleWithPermissions(id uint) (*Role, error) {
	var role Role
	err := r.db.Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ===== Permission operations =====

func (r *repositoryImpl) CreatePermission(permission *Permission) error {
	return r.db.Create(permission).Error
}

func (r *repositoryImpl) GetPermissionByID(id uint) (*Permission, error) {
	var permission Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *repositoryImpl) GetPermissionByName(name string) (*Permission, error) {
	var permission Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *repositoryImpl) UpdatePermission(permission *Permission) error {
	return r.db.Save(permission).Error
}

func (r *repositoryImpl) DeletePermission(id uint) error {
	return r.db.Delete(&Permission{}, id).Error
}

func (r *repositoryImpl) ListPermissions(query ListPermissionsQuery) ([]Permission, int64, error) {
	var permissions []Permission
	var total int64

	db := r.db.Model(&Permission{})

	// Apply filters
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR display_name ILIKE ?", searchPattern, searchPattern)
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	if query.Resource != "" {
		db = db.Where("resource = ?", query.Resource)
	}

	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}

	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	// Count total
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (query.Page - 1) * query.PageSize
	orderBy := fmt.Sprintf("%s %s", query.OrderBy, query.Order)

	err = db.Order(orderBy).Offset(offset).Limit(query.PageSize).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (r *repositoryImpl) GetPermissionsByIDs(ids []uint) ([]Permission, error) {
	var permissions []Permission
	err := r.db.Where("id IN ?", ids).Find(&permissions).Error
	return permissions, err
}

// ===== Role-Permission operations =====

func (r *repositoryImpl) AssignPermissionsToRole(roleID uint, permissionIDs []uint, grantedBy uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, remove existing permissions to handle de-selection
		if err := tx.Where("role_id = ?", roleID).Delete(&RolePermission{}).Error; err != nil {
			return err
		}

		if len(permissionIDs) == 0 {
			return nil // Nothing to assign
		}

		// Then, add the new set of permissions
		var rolePermissions []RolePermission
		for _, permissionID := range permissionIDs {
			rolePermissions = append(rolePermissions, RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			})
		}

		return tx.Create(&rolePermissions).Error
	})
}

func (r *repositoryImpl) RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error {
	return r.db.Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).Delete(&RolePermission{}).Error
}

func (r *repositoryImpl) GetRolePermissions(roleID uint) ([]Permission, error) {
	var permissions []Permission
	err := r.db.
		Joins("JOIN role_permissions on role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// ===== User-Role operations =====

func (r *repositoryImpl) AssignRoleToUser(userRole *UserRole) error {
	return r.db.Create(userRole).Error
}

func (r *repositoryImpl) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error
}

func (r *repositoryImpl) GetUserRoles(userID uint) ([]UserRole, error) {
	var userRoles []UserRole
	err := r.db.Where("user_id = ? AND is_active = true", userID).Preload("Role").Find(&userRoles).Error
	return userRoles, err
}

func (r *repositoryImpl) GetUsersWithRole(roleID uint) ([]UserRole, error) {
	var userRoles []UserRole
	err := r.db.Where("role_id = ?", roleID).Find(&userRoles).Error
	return userRoles, err
}

func (r *repositoryImpl) CheckUserRole(userID, roleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&UserRole{}).Where("user_id = ? AND role_id = ? AND is_active = true", userID, roleID).Count(&count).Error
	return count > 0, err
}

// ===== Organization-Role operations =====

func (r *repositoryImpl) AssignOrganizationRole(orgRole *OrganizationRole) error {
	return r.db.Create(orgRole).Error
}

func (r *repositoryImpl) RemoveOrganizationRole(userID, organizationID, roleID uint) error {
	return r.db.Where("user_id = ? AND organization_id = ? AND role_id = ?", userID, organizationID, roleID).
		Delete(&OrganizationRole{}).Error
}

func (r *repositoryImpl) GetUserOrganizationRoles(userID, organizationID uint) ([]OrganizationRole, error) {
	var orgRoles []OrganizationRole
	err := r.db.Where("user_id = ? AND organization_id = ? AND is_active = true",
		userID, organizationID).Preload("Role").Find(&orgRoles).Error
	return orgRoles, err
}

func (r *repositoryImpl) GetOrganizationUsers(organizationID, roleID uint) ([]OrganizationRole, error) {
	var orgRoles []OrganizationRole
	err := r.db.Where("organization_id = ? AND role_id = ? AND is_active = true",
		organizationID, roleID).Find(&orgRoles).Error
	return orgRoles, err
}

// ===== Team-Role operations =====

func (r *repositoryImpl) AssignTeamRole(teamRole *TeamRole) error {
	return r.db.Create(teamRole).Error
}

func (r *repositoryImpl) RemoveTeamRole(userID, teamID, roleID uint) error {
	return r.db.Where("user_id = ? AND team_id = ? AND role_id = ?", userID, teamID, roleID).
		Delete(&TeamRole{}).Error
}

func (r *repositoryImpl) GetUserTeamRoles(userID, teamID uint) ([]TeamRole, error) {
	var teamRoles []TeamRole
	err := r.db.Where("user_id = ? AND team_id = ? AND is_active = true",
		userID, teamID).Preload("Role").Find(&teamRoles).Error
	return teamRoles, err
}

func (r *repositoryImpl) GetTeamUsers(teamID, roleID uint) ([]TeamRole, error) {
	var teamRoles []TeamRole
	db := r.db.Where("team_id = ?", teamID)
	if roleID > 0 {
		db = db.Where("role_id = ?", roleID)
	}
	err := db.Find(&teamRoles).Error
	return teamRoles, err
}

// ===== Policy operations =====

func (r *repositoryImpl) CreatePolicy(policy *Policy) error {
	return r.db.Create(policy).Error
}

func (r *repositoryImpl) GetPolicyByID(id uint) (*Policy, error) {
	var policy Policy
	err := r.db.First(&policy, id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *repositoryImpl) UpdatePolicy(policy *Policy) error {
	return r.db.Save(policy).Error
}

func (r *repositoryImpl) DeletePolicy(id uint) error {
	return r.db.Delete(&Policy{}, id).Error
}

func (r *repositoryImpl) ListPolicies(query ListQuery) ([]Policy, int64, error) {
	var policies []Policy
	var total int64

	db := r.db.Model(&Policy{})

	// Apply filters from ListQuery
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("subject ILIKE ? OR action ILIKE ? OR object ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count total
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (query.Page - 1) * query.PageSize
	orderBy := fmt.Sprintf("%s %s", query.OrderBy, query.Order)

	err = db.Order(orderBy).Offset(offset).Limit(query.PageSize).Find(&policies).Error
	if err != nil {
		return nil, 0, err
	}

	return policies, total, nil
}

// ===== Permission checking operations =====

// GetUserAllPermissions gets all permissions for a user from their directly assigned roles
func (r *repositoryImpl) GetUserAllPermissions(userID uint) ([]string, error) {
	var permissions []string
	err := r.db.Raw(`
		SELECT DISTINCT p.name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND ur.is_active = true AND r.status = 1 AND p.status = 1
	`, userID).Scan(&permissions).Error
	return permissions, err
}

// GetUserOrganizationPermissions gets all permissions for a user within an organization
func (r *repositoryImpl) GetUserOrganizationPermissions(userID, organizationID uint) ([]string, error) {
	var permissions []string
	err := r.db.Raw(`
		SELECT DISTINCT p.name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		JOIN organization_roles org_r ON r.id = org_r.role_id
		WHERE org_r.user_id = ? AND org_r.organization_id = ?
		AND org_r.is_active = true AND r.status = 1 AND p.status = 1
	`, userID, organizationID).Scan(&permissions).Error
	return permissions, err
}

// GetUserTeamPermissions gets all permissions for a user within a team
func (r *repositoryImpl) GetUserTeamPermissions(userID, teamID uint) ([]string, error) {
	var permissions []string
	err := r.db.Raw(`
		SELECT DISTINCT p.name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		JOIN team_roles tr ON r.id = tr.role_id
		WHERE tr.user_id = ? AND tr.team_id = ?
		AND tr.is_active = true AND r.status = 1 AND p.status = 1
	`, userID, teamID).Scan(&permissions).Error
	return permissions, err
}

// CheckUserPermission checks if a user has a specific global permission
func (r *repositoryImpl) CheckUserPermission(userID uint, permission string) (bool, error) {
	var count int64

	err := r.db.Raw(`
		SELECT COUNT(DISTINCT p.id)
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND p.name = ?
		AND ur.is_active = true AND r.status = 1 AND p.status = 1
	`, userID, permission).Scan(&count).Error

	return count > 0, err
}

// CheckUserOrganizationPermission checks if a user has a specific permission in an organization
func (r *repositoryImpl) CheckUserOrganizationPermission(userID, organizationID uint, permission string) (bool, error) {
	var count int64

	err := r.db.Raw(`
		SELECT COUNT(DISTINCT p.id)
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		JOIN organization_roles org_r ON r.id = org_r.role_id
		WHERE org_r.user_id = ? AND org_r.organization_id = ? AND p.name = ?
		AND org_r.is_active = true AND r.status = 1 AND p.status = 1
	`, userID, organizationID, permission).Scan(&count).Error

	return count > 0, err
}

// CheckUserTeamPermission checks if a user has a specific permission in a team
func (r *repositoryImpl) CheckUserTeamPermission(userID, teamID uint, permission string) (bool, error) {
	permissions, err := r.GetUserTeamPermissions(userID, teamID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}
	return false, nil
}
