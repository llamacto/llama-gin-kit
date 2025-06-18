package organization

import (
	"context"
	"encoding/json"
	"errors"
)

// Role and Permission methods implementation

// CreateRole adds a new role
func (s *OrganizationServiceImpl) CreateRole(ctx context.Context, role *Role) error {
	// If organization ID is specified, verify organization exists
	if role.OrganizationID != nil && *role.OrganizationID > 0 {
		_, err := s.GetOrganization(ctx, *role.OrganizationID)
		if err != nil {
			return errors.New("organization not found")
		}
	}
	
	return s.repo.CreateRole(ctx, role)
}

// UpdateRole updates an existing role
func (s *OrganizationServiceImpl) UpdateRole(ctx context.Context, role *Role) error {
	// Check if role exists
	existingRole, err := s.GetRole(ctx, role.ID)
	if err != nil {
		return errors.New("role not found")
	}
	
	// Prevent change of organization
	if (role.OrganizationID == nil && existingRole.OrganizationID != nil) ||
	   (role.OrganizationID != nil && existingRole.OrganizationID == nil) ||
	   (role.OrganizationID != nil && existingRole.OrganizationID != nil && 
	    *role.OrganizationID != *existingRole.OrganizationID) {
		return errors.New("cannot change role's organization")
	}
	
	return s.repo.UpdateRole(ctx, role)
}

// DeleteRole removes a role by ID
func (s *OrganizationServiceImpl) DeleteRole(ctx context.Context, id uint) error {
	// Check if role exists
	role, err := s.GetRole(ctx, id)
	if err != nil {
		return errors.New("role not found")
	}
	
	// Check if it's the default role
	if role.IsDefault {
		return errors.New("cannot delete default role")
	}
	
	// Check if role is in use
	var count int64
	if err := s.db.Model(&Member{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("role is in use and cannot be deleted")
	}
	
	return s.repo.DeleteRole(ctx, id)
}

// GetRole retrieves a role by ID
func (s *OrganizationServiceImpl) GetRole(ctx context.Context, id uint) (*Role, error) {
	return s.repo.GetRole(ctx, id)
}

// ListRoles retrieves roles for an organization with pagination
func (s *OrganizationServiceImpl) ListRoles(ctx context.Context, orgID uint, page, pageSize int) ([]*Role, int64, error) {
	// If organization ID is specified, verify organization exists
	if orgID > 0 {
		_, err := s.GetOrganization(ctx, orgID)
		if err != nil {
			return nil, 0, errors.New("organization not found")
		}
	}
	
	return s.repo.ListRoles(ctx, orgID, page, pageSize)
}

// CheckPermission checks if a user has a specific permission in an organization
func (s *OrganizationServiceImpl) CheckPermission(ctx context.Context, userID uint, orgID uint, permission string) (bool, error) {
	// Get member record
	member, err := s.repo.GetMemberByUserAndOrg(ctx, userID, orgID)
	if err != nil {
		return false, errors.New("user is not a member of this organization")
	}
	
	// Check if member is active
	if member.Status != 1 {
		return false, nil
	}
	
	// Get role
	role, err := s.GetRole(ctx, member.RoleID)
	if err != nil {
		return false, errors.New("member role not found")
	}
	
	// Parse permissions
	var permissions map[string]interface{}
	if err := json.Unmarshal([]byte(role.Permissions), &permissions); err != nil {
		return false, errors.New("invalid permission format")
	}
	
	// Check wildcard permission
	if val, ok := permissions["*"]; ok {
		if boolVal, ok := val.(bool); ok && boolVal {
			return true, nil
		}
	}
	
	// Check specific permission
	if val, ok := permissions[permission]; ok {
		if boolVal, ok := val.(bool); ok && boolVal {
			return true, nil
		}
	}
	
	return false, nil
}
