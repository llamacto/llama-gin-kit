package authorization

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Service interface for authorization business logic
type Service interface {
	// Role management
	CreateRole(req CreateRoleRequest, createdBy uint) (*RoleResponse, error)
	GetRole(id uint) (*RoleResponse, error)
	UpdateRole(id uint, req UpdateRoleRequest, updatedBy uint) (*RoleResponse, error)
	DeleteRole(id uint, deletedBy uint) error
	ListRoles(query ListRolesQuery) (*ListResponse, error)
	GetRoleWithPermissions(id uint) (*RoleWithPermissionsResponse, error)

	// Permission management
	CreatePermission(req CreatePermissionRequest, createdBy uint) (*PermissionResponse, error)
	GetPermission(id uint) (*PermissionResponse, error)
	UpdatePermission(id uint, req UpdatePermissionRequest, updatedBy uint) (*PermissionResponse, error)
	DeletePermission(id uint, deletedBy uint) error
	ListPermissions(query ListPermissionsQuery) (*ListResponse, error)

	// Role-Permission management
	AssignPermissionsToRole(roleID uint, req AssignPermissionsRequest, assignedBy uint) error
	RemovePermissionsFromRole(roleID uint, req RemovePermissionsRequest, removedBy uint) error

	// User-Role management
	AssignRoleToUser(req AssignRoleRequest, assignedBy uint) (*UserRoleResponse, error)
	AssignRolesToUser(req AssignRolesRequest, assignedBy uint) ([]UserRoleResponse, error)
	RemoveRoleFromUser(userID, roleID uint, removedBy uint) error
	GetUserRoles(userID uint) ([]UserRoleResponse, error)

	// Organization-Role management
	AssignOrganizationRole(req AssignOrganizationRoleRequest, assignedBy uint) (*OrganizationRoleResponse, error)
	RemoveOrganizationRole(userID, organizationID, roleID uint, removedBy uint) error
	GetUserOrganizationRoles(userID, organizationID uint) ([]OrganizationRoleResponse, error)

	// Team-Role management
	AssignTeamRole(req AssignTeamRoleRequest, assignedBy uint) (*TeamRoleResponse, error)
	RemoveTeamRole(userID, teamID, roleID uint, removedBy uint) error
	GetUserTeamRoles(userID, teamID uint) ([]TeamRoleResponse, error)

	// Permission checking
	CheckPermission(req CheckPermissionRequest) (*CheckPermissionResponse, error)
	GetUserPermissionsSummary(userID uint) (*UserPermissionsSummaryResponse, error)
	HasPermission(userID uint, permission string) (bool, error)
	HasOrganizationPermission(userID, organizationID uint, permission string) (bool, error)
	HasTeamPermission(userID, teamID uint, permission string) (bool, error)

	// System initialization
	InitializeSystemRoles() error
	InitializeSystemPermissions() error
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	repo Repository
}

// NewService creates a new authorization service
func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

// ===== Role management =====

func (s *serviceImpl) CreateRole(req CreateRoleRequest, createdBy uint) (*RoleResponse, error) {
	// Check if role name already exists
	existingRole, err := s.repo.GetRoleByName(req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existingRole != nil {
		return nil, errors.New("role name already exists")
	}

	// Create role
	role := &Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Level:       req.Level,
		Status:      req.Status,
		IsSystem:    false,
	}

	err = s.repo.CreateRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return s.roleToResponse(role), nil
}

func (s *serviceImpl) GetRole(id uint) (*RoleResponse, error) {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return s.roleToResponse(role), nil
}

func (s *serviceImpl) UpdateRole(id uint, req UpdateRoleRequest, updatedBy uint) (*RoleResponse, error) {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Check if it's a system role
	if role.IsSystem {
		return nil, errors.New("cannot update system role")
	}

	// Update fields
	if req.DisplayName != nil {
		role.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Level != nil {
		role.Level = *req.Level
	}
	if req.Status != nil {
		role.Status = *req.Status
	}

	err = s.repo.UpdateRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return s.roleToResponse(role), nil
}

func (s *serviceImpl) DeleteRole(id uint, deletedBy uint) error {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Check if it's a system role
	if role.IsSystem {
		return errors.New("cannot delete system role")
	}

	// Check if role is assigned to users
	users, err := s.repo.GetUsersWithRole(id)
	if err != nil {
		return fmt.Errorf("failed to check role assignments: %w", err)
	}
	if len(users) > 0 {
		return errors.New("cannot delete role that is assigned to users")
	}

	return s.repo.DeleteRole(id)
}

func (s *serviceImpl) ListRoles(query ListRolesQuery) (*ListResponse, error) {
	roles, total, err := s.repo.ListRoles(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	var responses []RoleResponse
	for _, role := range roles {
		responses = append(responses, *s.roleToResponse(&role))
	}

	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

	return &ListResponse{
		Data:       responses,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *serviceImpl) GetRoleWithPermissions(id uint) (*RoleWithPermissionsResponse, error) {
	role, err := s.repo.GetRoleWithPermissions(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role with permissions: %w", err)
	}

	var permissions []PermissionResponse
	for _, perm := range role.Permissions {
		permissions = append(permissions, *s.permissionToResponse(&perm))
	}

	return &RoleWithPermissionsResponse{
		RoleResponse: *s.roleToResponse(role),
		Permissions:  permissions,
	}, nil
}

// ===== Permission management =====

func (s *serviceImpl) CreatePermission(req CreatePermissionRequest, createdBy uint) (*PermissionResponse, error) {
	// Check if permission name already exists
	existingPerm, err := s.repo.GetPermissionByName(req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing permission: %w", err)
	}
	if existingPerm != nil {
		return nil, errors.New("permission name already exists")
	}

	// Create permission
	permission := &Permission{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
		Status:      req.Status,
		IsSystem:    false,
	}

	err = s.repo.CreatePermission(permission)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return s.permissionToResponse(permission), nil
}

func (s *serviceImpl) GetPermission(id uint) (*PermissionResponse, error) {
	permission, err := s.repo.GetPermissionByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return s.permissionToResponse(permission), nil
}

func (s *serviceImpl) UpdatePermission(id uint, req UpdatePermissionRequest, updatedBy uint) (*PermissionResponse, error) {
	permission, err := s.repo.GetPermissionByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	// Check if it's a system permission
	if permission.IsSystem {
		return nil, errors.New("cannot update system permission")
	}

	// Update fields
	if req.DisplayName != nil {
		permission.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}
	if req.Resource != nil {
		permission.Resource = *req.Resource
	}
	if req.Action != nil {
		permission.Action = *req.Action
	}
	if req.Category != nil {
		permission.Category = *req.Category
	}
	if req.Status != nil {
		permission.Status = *req.Status
	}

	err = s.repo.UpdatePermission(permission)
	if err != nil {
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	return s.permissionToResponse(permission), nil
}

func (s *serviceImpl) DeletePermission(id uint, deletedBy uint) error {
	permission, err := s.repo.GetPermissionByID(id)
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}

	// Check if it's a system permission
	if permission.IsSystem {
		return errors.New("cannot delete system permission")
	}

	return s.repo.DeletePermission(id)
}

func (s *serviceImpl) ListPermissions(query ListPermissionsQuery) (*ListResponse, error) {
	permissions, total, err := s.repo.ListPermissions(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	var responses []PermissionResponse
	for _, perm := range permissions {
		responses = append(responses, *s.permissionToResponse(&perm))
	}

	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

	return &ListResponse{
		Data:       responses,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ===== Role-Permission management =====

func (s *serviceImpl) AssignPermissionsToRole(roleID uint, req AssignPermissionsRequest, assignedBy uint) error {
	// Verify role exists
	_, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Verify permissions exist
	permissions, err := s.repo.GetPermissionsByIDs(req.PermissionIDs)
	if err != nil {
		return fmt.Errorf("failed to get permissions: %w", err)
	}
	if len(permissions) != len(req.PermissionIDs) {
		return errors.New("some permissions not found")
	}

	return s.repo.AssignPermissionsToRole(roleID, req.PermissionIDs, assignedBy)
}

func (s *serviceImpl) RemovePermissionsFromRole(roleID uint, req RemovePermissionsRequest, removedBy uint) error {
	// Verify role exists
	_, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	return s.repo.RemovePermissionsFromRole(roleID, req.PermissionIDs)
}

// ===== User-Role management =====

func (s *serviceImpl) AssignRoleToUser(req AssignRoleRequest, assignedBy uint) (*UserRoleResponse, error) {
	// Verify role exists
	role, err := s.repo.GetRoleByID(req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Check if user already has this role
	hasRole, err := s.repo.CheckUserRole(req.UserID, req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user role: %w", err)
	}
	if hasRole {
		return nil, errors.New("user already has this role")
	}

	userRole := &UserRole{
		UserID:     req.UserID,
		RoleID:     req.RoleID,
		AssignedBy: assignedBy,
		ExpiresAt:  req.ExpiresAt,
		IsActive:   true,
	}

	err = s.repo.AssignRoleToUser(userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to assign role to user: %w", err)
	}

	return &UserRoleResponse{
		ID:         userRole.ID,
		UserID:     userRole.UserID,
		RoleID:     userRole.RoleID,
		AssignedBy: userRole.AssignedBy,
		ExpiresAt:  userRole.ExpiresAt,
		IsActive:   userRole.IsActive,
		CreatedAt:  userRole.CreatedAt,
		UpdatedAt:  userRole.UpdatedAt,
		Role:       *s.roleToResponse(role),
	}, nil
}

func (s *serviceImpl) AssignRolesToUser(req AssignRolesRequest, assignedBy uint) ([]UserRoleResponse, error) {
	var responses []UserRoleResponse

	for _, roleID := range req.RoleIDs {
		assignReq := AssignRoleRequest{
			UserID: req.UserID,
			RoleID: roleID,
		}

		response, err := s.AssignRoleToUser(assignReq, assignedBy)
		if err != nil {
			// Continue with other roles, but log the error
			continue
		}

		responses = append(responses, *response)
	}

	return responses, nil
}

func (s *serviceImpl) RemoveRoleFromUser(userID, roleID uint, removedBy uint) error {
	return s.repo.RemoveRoleFromUser(userID, roleID)
}

func (s *serviceImpl) GetUserRoles(userID uint) ([]UserRoleResponse, error) {
	userRoles, err := s.repo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	var responses []UserRoleResponse
	for _, userRole := range userRoles {
		responses = append(responses, UserRoleResponse{
			ID:         userRole.ID,
			UserID:     userRole.UserID,
			RoleID:     userRole.RoleID,
			AssignedBy: userRole.AssignedBy,
			ExpiresAt:  userRole.ExpiresAt,
			IsActive:   userRole.IsActive,
			CreatedAt:  userRole.CreatedAt,
			UpdatedAt:  userRole.UpdatedAt,
			Role:       *s.roleToResponse(&userRole.Role),
		})
	}

	return responses, nil
}

// ===== Organization-Role management =====

func (s *serviceImpl) AssignOrganizationRole(req AssignOrganizationRoleRequest, assignedBy uint) (*OrganizationRoleResponse, error) {
	// Verify role exists
	role, err := s.repo.GetRoleByID(req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	orgRole := &OrganizationRole{
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		RoleID:         req.RoleID,
		AssignedBy:     assignedBy,
		IsActive:       true,
	}

	err = s.repo.AssignOrganizationRole(orgRole)
	if err != nil {
		return nil, fmt.Errorf("failed to assign organization role: %w", err)
	}

	return &OrganizationRoleResponse{
		ID:             orgRole.ID,
		UserID:         orgRole.UserID,
		OrganizationID: orgRole.OrganizationID,
		RoleID:         orgRole.RoleID,
		AssignedBy:     orgRole.AssignedBy,
		IsActive:       orgRole.IsActive,
		CreatedAt:      orgRole.CreatedAt,
		UpdatedAt:      orgRole.UpdatedAt,
		Role:           *s.roleToResponse(role),
	}, nil
}

func (s *serviceImpl) RemoveOrganizationRole(userID, organizationID, roleID uint, removedBy uint) error {
	return s.repo.RemoveOrganizationRole(userID, organizationID, roleID)
}

func (s *serviceImpl) GetUserOrganizationRoles(userID, organizationID uint) ([]OrganizationRoleResponse, error) {
	orgRoles, err := s.repo.GetUserOrganizationRoles(userID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organization roles: %w", err)
	}

	var responses []OrganizationRoleResponse
	for _, orgRole := range orgRoles {
		responses = append(responses, OrganizationRoleResponse{
			ID:             orgRole.ID,
			UserID:         orgRole.UserID,
			OrganizationID: orgRole.OrganizationID,
			RoleID:         orgRole.RoleID,
			AssignedBy:     orgRole.AssignedBy,
			IsActive:       orgRole.IsActive,
			CreatedAt:      orgRole.CreatedAt,
			UpdatedAt:      orgRole.UpdatedAt,
			Role:           *s.roleToResponse(&orgRole.Role),
		})
	}

	return responses, nil
}

// ===== Team-Role management =====

func (s *serviceImpl) AssignTeamRole(req AssignTeamRoleRequest, assignedBy uint) (*TeamRoleResponse, error) {
	// Verify role exists
	role, err := s.repo.GetRoleByID(req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	teamRole := &TeamRole{
		UserID:     req.UserID,
		TeamID:     req.TeamID,
		RoleID:     req.RoleID,
		AssignedBy: assignedBy,
		IsActive:   true,
	}

	err = s.repo.AssignTeamRole(teamRole)
	if err != nil {
		return nil, fmt.Errorf("failed to assign team role: %w", err)
	}

	return &TeamRoleResponse{
		ID:         teamRole.ID,
		UserID:     teamRole.UserID,
		TeamID:     teamRole.TeamID,
		RoleID:     teamRole.RoleID,
		AssignedBy: teamRole.AssignedBy,
		IsActive:   teamRole.IsActive,
		CreatedAt:  teamRole.CreatedAt,
		UpdatedAt:  teamRole.UpdatedAt,
		Role:       *s.roleToResponse(role),
	}, nil
}

func (s *serviceImpl) RemoveTeamRole(userID, teamID, roleID uint, removedBy uint) error {
	return s.repo.RemoveTeamRole(userID, teamID, roleID)
}

func (s *serviceImpl) GetUserTeamRoles(userID, teamID uint) ([]TeamRoleResponse, error) {
	teamRoles, err := s.repo.GetUserTeamRoles(userID, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user team roles: %w", err)
	}

	var responses []TeamRoleResponse
	for _, teamRole := range teamRoles {
		responses = append(responses, TeamRoleResponse{
			ID:         teamRole.ID,
			UserID:     teamRole.UserID,
			TeamID:     teamRole.TeamID,
			RoleID:     teamRole.RoleID,
			AssignedBy: teamRole.AssignedBy,
			IsActive:   teamRole.IsActive,
			CreatedAt:  teamRole.CreatedAt,
			UpdatedAt:  teamRole.UpdatedAt,
			Role:       *s.roleToResponse(&teamRole.Role),
		})
	}

	return responses, nil
}

// ===== Permission checking =====

func (s *serviceImpl) CheckPermission(req CheckPermissionRequest) (*CheckPermissionResponse, error) {
	var hasPermission bool
	var roles []string
	var source string

	// Check global permissions first
	globalHasPermission, err := s.repo.CheckUserPermission(req.UserID, req.Permission)
	if err != nil {
		return nil, fmt.Errorf("failed to check global permission: %w", err)
	}

	if globalHasPermission {
		hasPermission = true
		source = "global"

		// Get user's global roles
		userRoles, err := s.repo.GetUserRoles(req.UserID)
		if err == nil {
			for _, userRole := range userRoles {
				roles = append(roles, userRole.Role.Name)
			}
		}
	}

	// Check organization permissions if organization ID is provided
	if !hasPermission && req.OrganizationID != nil {
		orgHasPermission, err := s.repo.CheckUserOrganizationPermission(req.UserID, *req.OrganizationID, req.Permission)
		if err != nil {
			return nil, fmt.Errorf("failed to check organization permission: %w", err)
		}

		if orgHasPermission {
			hasPermission = true
			source = "organization"

			// Get user's organization roles
			orgRoles, err := s.repo.GetUserOrganizationRoles(req.UserID, *req.OrganizationID)
			if err == nil {
				for _, orgRole := range orgRoles {
					roles = append(roles, orgRole.Role.Name)
				}
			}
		}
	}

	// Check team permissions if team ID is provided
	if !hasPermission && req.TeamID != nil {
		teamHasPermission, err := s.repo.CheckUserTeamPermission(req.UserID, *req.TeamID, req.Permission)
		if err != nil {
			return nil, fmt.Errorf("failed to check team permission: %w", err)
		}

		if teamHasPermission {
			hasPermission = true
			source = "team"

			// Get user's team roles
			teamRoles, err := s.repo.GetUserTeamRoles(req.UserID, *req.TeamID)
			if err == nil {
				for _, teamRole := range teamRoles {
					roles = append(roles, teamRole.Role.Name)
				}
			}
		}
	}

	return &CheckPermissionResponse{
		HasPermission: hasPermission,
		UserID:        req.UserID,
		Permission:    req.Permission,
		Resource:      req.Resource,
		Roles:         roles,
		Source:        source,
	}, nil
}

func (s *serviceImpl) GetUserPermissionsSummary(userID uint) (*UserPermissionsSummaryResponse, error) {
	// Get global roles
	globalRoles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Get all permissions
	allPermissions, err := s.repo.GetUserAllPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Convert to role responses
	var globalRoleResponses []RoleResponse
	for _, userRole := range globalRoles {
		globalRoleResponses = append(globalRoleResponses, userRole.Role)
	}

	// Get effective permissions (unique)
	permissionMap := make(map[string]bool)
	for _, perm := range allPermissions {
		permissionMap[perm] = true
	}

	var effectivePermissions []PermissionResponse
	for permName := range permissionMap {
		perm, err := s.repo.GetPermissionByName(permName)
		if err == nil {
			effectivePermissions = append(effectivePermissions, *s.permissionToResponse(perm))
		}
	}

	return &UserPermissionsSummaryResponse{
		UserID:               userID,
		GlobalRoles:          globalRoleResponses,
		OrganizationRoles:    []OrganizationRoleResponse{}, // TODO: Get organization roles
		TeamRoles:            []TeamRoleResponse{},         // TODO: Get team roles
		AllPermissions:       allPermissions,
		EffectivePermissions: effectivePermissions,
	}, nil
}

func (s *serviceImpl) HasPermission(userID uint, permission string) (bool, error) {
	return s.repo.CheckUserPermission(userID, permission)
}

func (s *serviceImpl) HasOrganizationPermission(userID, organizationID uint, permission string) (bool, error) {
	return s.repo.CheckUserOrganizationPermission(userID, organizationID, permission)
}

func (s *serviceImpl) HasTeamPermission(userID, teamID uint, permission string) (bool, error) {
	return s.repo.CheckUserTeamPermission(userID, teamID, permission)
}

// ===== System initialization =====

func (s *serviceImpl) InitializeSystemRoles() error {
	systemRoles := []Role{
		{
			Name:        "super_admin",
			DisplayName: "Super Administrator",
			Description: "Super administrator with all permissions",
			Level:       1000,
			IsSystem:    true,
			Status:      1,
		},
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "System administrator",
			Level:       900,
			IsSystem:    true,
			Status:      1,
		},
		{
			Name:        "moderator",
			DisplayName: "Moderator",
			Description: "Content moderator",
			Level:       500,
			IsSystem:    true,
			Status:      1,
		},
		{
			Name:        "user",
			DisplayName: "Regular User",
			Description: "Regular user with basic permissions",
			Level:       100,
			IsSystem:    true,
			Status:      1,
		},
	}

	for _, role := range systemRoles {
		// Check if role already exists
		existingRole, err := s.repo.GetRoleByName(role.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing role: %w", err)
		}

		if existingRole == nil {
			err = s.repo.CreateRole(&role)
			if err != nil {
				return fmt.Errorf("failed to create system role %s: %w", role.Name, err)
			}
		}
	}

	return nil
}

func (s *serviceImpl) InitializeSystemPermissions() error {
	systemPermissions := []Permission{
		// User management
		{Name: "users.create", DisplayName: "Create Users", Resource: "users", Action: "create", Category: "user_management", IsSystem: true, Status: 1},
		{Name: "users.read", DisplayName: "Read Users", Resource: "users", Action: "read", Category: "user_management", IsSystem: true, Status: 1},
		{Name: "users.update", DisplayName: "Update Users", Resource: "users", Action: "update", Category: "user_management", IsSystem: true, Status: 1},
		{Name: "users.delete", DisplayName: "Delete Users", Resource: "users", Action: "delete", Category: "user_management", IsSystem: true, Status: 1},

		// Organization management
		{Name: "organizations.create", DisplayName: "Create Organizations", Resource: "organizations", Action: "create", Category: "organization_management", IsSystem: true, Status: 1},
		{Name: "organizations.read", DisplayName: "Read Organizations", Resource: "organizations", Action: "read", Category: "organization_management", IsSystem: true, Status: 1},
		{Name: "organizations.update", DisplayName: "Update Organizations", Resource: "organizations", Action: "update", Category: "organization_management", IsSystem: true, Status: 1},
		{Name: "organizations.delete", DisplayName: "Delete Organizations", Resource: "organizations", Action: "delete", Category: "organization_management", IsSystem: true, Status: 1},

		// Team management
		{Name: "teams.create", DisplayName: "Create Teams", Resource: "teams", Action: "create", Category: "team_management", IsSystem: true, Status: 1},
		{Name: "teams.read", DisplayName: "Read Teams", Resource: "teams", Action: "read", Category: "team_management", IsSystem: true, Status: 1},
		{Name: "teams.update", DisplayName: "Update Teams", Resource: "teams", Action: "update", Category: "team_management", IsSystem: true, Status: 1},
		{Name: "teams.delete", DisplayName: "Delete Teams", Resource: "teams", Action: "delete", Category: "team_management", IsSystem: true, Status: 1},

		// Role management
		{Name: "roles.create", DisplayName: "Create Roles", Resource: "roles", Action: "create", Category: "role_management", IsSystem: true, Status: 1},
		{Name: "roles.read", DisplayName: "Read Roles", Resource: "roles", Action: "read", Category: "role_management", IsSystem: true, Status: 1},
		{Name: "roles.update", DisplayName: "Update Roles", Resource: "roles", Action: "update", Category: "role_management", IsSystem: true, Status: 1},
		{Name: "roles.delete", DisplayName: "Delete Roles", Resource: "roles", Action: "delete", Category: "role_management", IsSystem: true, Status: 1},

		// Permission management
		{Name: "permissions.create", DisplayName: "Create Permissions", Resource: "permissions", Action: "create", Category: "permission_management", IsSystem: true, Status: 1},
		{Name: "permissions.read", DisplayName: "Read Permissions", Resource: "permissions", Action: "read", Category: "permission_management", IsSystem: true, Status: 1},
		{Name: "permissions.update", DisplayName: "Update Permissions", Resource: "permissions", Action: "update", Category: "permission_management", IsSystem: true, Status: 1},
		{Name: "permissions.delete", DisplayName: "Delete Permissions", Resource: "permissions", Action: "delete", Category: "permission_management", IsSystem: true, Status: 1},
	}

	for _, permission := range systemPermissions {
		// Check if permission already exists
		existingPerm, err := s.repo.GetPermissionByName(permission.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing permission: %w", err)
		}

		if existingPerm == nil {
			err = s.repo.CreatePermission(&permission)
			if err != nil {
				return fmt.Errorf("failed to create system permission %s: %w", permission.Name, err)
			}
		}
	}

	return nil
}

// ===== Helper methods =====

func (s *serviceImpl) roleToResponse(role *Role) *RoleResponse {
	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		Level:       role.Level,
		IsSystem:    role.IsSystem,
		Status:      role.Status,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func (s *serviceImpl) permissionToResponse(permission *Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		DisplayName: permission.DisplayName,
		Description: permission.Description,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Category:    permission.Category,
		IsSystem:    permission.IsSystem,
		Status:      permission.Status,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}
