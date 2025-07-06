// File: internal/domain/services/permission.go
// Tạo tại: internal/domain/services/permission.go
// Mục đích: Service quản lý phân quyền chi tiết

package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/dto/request"
	"github.com/godiidev/appsynex/internal/dto/response"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
)

type PermissionService interface {
	// Permission CRUD
	GetAllPermissions() (*response.PermissionsResponse, error)
	GetPermissionsByModule(module string) (*response.PermissionsResponse, error)
	GetPermissionGroups() (*response.PermissionGroupsResponse, error)
	CreatePermission(req request.CreatePermissionRequest) (*response.PermissionResponse, error)
	UpdatePermission(id uint, req request.UpdatePermissionRequest) (*response.PermissionResponse, error)
	DeletePermission(id uint) error
	
	// Role Permission Management
	AssignPermissionsToRole(roleID uint, permissionIDs []uint, grantedBy uint) error
	RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) (*response.RolePermissionsResponse, error)
	
	// User Permission Management (direct permissions)
	GrantUserPermission(req request.GrantUserPermissionRequest) error
	RevokeUserPermission(userID uint, permissionID uint) error
	GetUserPermissions(userID uint) (*response.UserPermissionsResponse, error)
	
	// Permission Checking
	UserHasPermission(userID uint, module string, action string, resource ...string) (bool, error)
	GetUserEffectivePermissions(userID uint) (*response.EffectivePermissionsResponse, error)
	
	// Bulk Operations
	BulkAssignPermissions(req request.BulkAssignPermissionsRequest) error
	CloneRolePermissions(fromRoleID, toRoleID uint, grantedBy uint) error
}

type permissionService struct {
	permissionRepo interfaces.PermissionRepository
	roleRepo       interfaces.RoleRepository
	userRepo       interfaces.UserRepository
}

func NewPermissionService(
	permissionRepo interfaces.PermissionRepository,
	roleRepo interfaces.RoleRepository,
	userRepo interfaces.UserRepository,
) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		userRepo:       userRepo,
	}
}

func (s *permissionService) GetAllPermissions() (*response.PermissionsResponse, error) {
	permissions, err := s.permissionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Group by module
	moduleGroups := make(map[string][]response.PermissionResponse)
	
	for _, perm := range permissions {
		if _, exists := moduleGroups[perm.Module]; !exists {
			moduleGroups[perm.Module] = []response.PermissionResponse{}
		}
		
		moduleGroups[perm.Module] = append(moduleGroups[perm.Module], response.PermissionResponse{
			ID:             perm.ID,
			Module:         perm.Module,
			Action:         perm.Action,
			Resource:       perm.Resource,
			PermissionName: perm.PermissionName,
			Description:    perm.Description,
			IsActive:       perm.IsActive,
		})
	}

	return &response.PermissionsResponse{
		Permissions:  convertPermissionsToResponse(permissions),
		ModuleGroups: moduleGroups,
		Total:        len(permissions),
	}, nil
}

func (s *permissionService) GetPermissionsByModule(module string) (*response.PermissionsResponse, error) {
	permissions, err := s.permissionRepo.FindByModule(module)
	if err != nil {
		return nil, err
	}

	return &response.PermissionsResponse{
		Permissions: convertPermissionsToResponse(permissions),
		Total:       len(permissions),
	}, nil
}

func (s *permissionService) GetPermissionGroups() (*response.PermissionGroupsResponse, error) {
	groups, err := s.permissionRepo.FindAllGroups()
	if err != nil {
		return nil, err
	}

	var groupResponses []response.PermissionGroupResponse
	for _, group := range groups {
		// Get permissions for this group
		permissions, _ := s.permissionRepo.FindByModule(group.Module)
		
		groupResponses = append(groupResponses, response.PermissionGroupResponse{
			ID:          group.ID,
			GroupName:   group.GroupName,
			DisplayName: group.DisplayName,
			Description: group.Description,
			Module:      group.Module,
			SortOrder:   group.SortOrder,
			IsActive:    group.IsActive,
			Permissions: convertPermissionsToResponse(permissions),
		})
	}

	return &response.PermissionGroupsResponse{
		Groups: groupResponses,
		Total:  len(groupResponses),
	}, nil
}

func (s *permissionService) CreatePermission(req request.CreatePermissionRequest) (*response.PermissionResponse, error) {
	// Generate permission name if not provided
	permissionName := req.PermissionName
	if permissionName == "" {
		if req.Resource != "" {
			permissionName = fmt.Sprintf("%s_%s_%s", req.Module, req.Action, req.Resource)
		} else {
			permissionName = fmt.Sprintf("%s_%s", req.Module, req.Action)
		}
	}

	// Check if permission already exists
	existing, _ := s.permissionRepo.FindByName(permissionName)
	if existing != nil {
		return nil, errors.New("permission already exists")
	}

	permission := &models.Permission{
		Module:         strings.ToUpper(req.Module),
		Action:         strings.ToUpper(req.Action),
		Resource:       req.Resource,
		PermissionName: strings.ToUpper(permissionName),
		Description:    req.Description,
		IsActive:       true,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return &response.PermissionResponse{
		ID:             permission.ID,
		Module:         permission.Module,
		Action:         permission.Action,
		Resource:       permission.Resource,
		PermissionName: permission.PermissionName,
		Description:    permission.Description,
		IsActive:       permission.IsActive,
	}, nil
}

func (s *permissionService) UpdatePermission(id uint, req request.UpdatePermissionRequest) (*response.PermissionResponse, error) {
	permission, err := s.permissionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("permission not found")
	}

	// Update fields if provided
	if req.Description != "" {
		permission.Description = req.Description
	}
	if req.IsActive != nil {
		permission.IsActive = *req.IsActive
	}

	if err := s.permissionRepo.Update(permission); err != nil {
		return nil, err
	}

	return &response.PermissionResponse{
		ID:             permission.ID,
		Module:         permission.Module,
		Action:         permission.Action,
		Resource:       permission.Resource,
		PermissionName: permission.PermissionName,
		Description:    permission.Description,
		IsActive:       permission.IsActive,
	}, nil
}

func (s *permissionService) DeletePermission(id uint) error {
	// Check if permission exists
	_, err := s.permissionRepo.FindByID(id)
	if err != nil {
		return errors.New("permission not found")
	}

	// Check if permission is in use
	inUse, err := s.permissionRepo.IsPermissionInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return errors.New("cannot delete permission that is currently assigned")
	}

	return s.permissionRepo.Delete(id)
}

func (s *permissionService) AssignPermissionsToRole(roleID uint, permissionIDs []uint, grantedBy uint) error {
	// Validate role exists
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Validate permissions exist
	for _, permID := range permissionIDs {
		_, err := s.permissionRepo.FindByID(permID)
		if err != nil {
			return fmt.Errorf("permission %d not found", permID)
		}
	}

	return s.permissionRepo.AssignPermissionsToRole(roleID, permissionIDs, grantedBy)
}

func (s *permissionService) RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error {
	// Validate role exists
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	return s.permissionRepo.RemovePermissionsFromRole(roleID, permissionIDs)
}

func (s *permissionService) GetRolePermissions(roleID uint) (*response.RolePermissionsResponse, error) {
	role, err := s.roleRepo.FindByIDWithPermissions(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Group permissions by module
	moduleGroups := make(map[string][]response.PermissionResponse)
	for _, perm := range role.Permissions {
		if _, exists := moduleGroups[perm.Module]; !exists {
			moduleGroups[perm.Module] = []response.PermissionResponse{}
		}
		
		moduleGroups[perm.Module] = append(moduleGroups[perm.Module], response.PermissionResponse{
			ID:             perm.ID,
			Module:         perm.Module,
			Action:         perm.Action,
			Resource:       perm.Resource,
			PermissionName: perm.PermissionName,
			Description:    perm.Description,
			IsActive:       perm.IsActive,
		})
	}

	return &response.RolePermissionsResponse{
		RoleID:       roleID,
		RoleName:     role.RoleName,
		Permissions:  convertPermissionsToResponse(role.Permissions),
		ModuleGroups: moduleGroups,
		Total:        len(role.Permissions),
	}, nil
}

func (s *permissionService) GrantUserPermission(req request.GrantUserPermissionRequest) error {
	// Validate user exists
	_, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Validate permission exists
	_, err = s.permissionRepo.FindByID(req.PermissionID)
	if err != nil {
		return errors.New("permission not found")
	}

	return s.permissionRepo.GrantUserPermission(req)
}

func (s *permissionService) RevokeUserPermission(userID uint, permissionID uint) error {
	return s.permissionRepo.RevokeUserPermission(userID, permissionID)
}

func (s *permissionService) GetUserPermissions(userID uint) (*response.UserPermissionsResponse, error) {
	user, err := s.userRepo.FindByIDWithRoles(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get role-based permissions
	rolePermissions, err := s.permissionRepo.GetUserRolePermissions(userID)
	if err != nil {
		return nil, err
	}

	// Get direct user permissions
	directPermissions, err := s.permissionRepo.GetUserDirectPermissions(userID)
	if err != nil {
		return nil, err
	}

	return &response.UserPermissionsResponse{
		UserID:            userID,
		Username:          user.Username,
		RolePermissions:   convertPermissionsToResponse(rolePermissions),
		DirectPermissions: convertUserPermissionsToResponse(directPermissions),
	}, nil
}

func (s *permissionService) UserHasPermission(userID uint, module string, action string, resource ...string) (bool, error) {
	resourceStr := ""
	if len(resource) > 0 {
		resourceStr = resource[0]
	}

	return s.permissionRepo.UserHasPermission(userID, module, action, resourceStr)
}

func (s *permissionService) GetUserEffectivePermissions(userID uint) (*response.EffectivePermissionsResponse, error) {
	effectivePermissions, err := s.permissionRepo.GetUserEffectivePermissions(userID)
	if err != nil {
		return nil, err
	}

	// Group by module
	moduleGroups := make(map[string][]response.PermissionResponse)
	for _, perm := range effectivePermissions {
		if _, exists := moduleGroups[perm.Module]; !exists {
			moduleGroups[perm.Module] = []response.PermissionResponse{}
		}
		
		moduleGroups[perm.Module] = append(moduleGroups[perm.Module], response.PermissionResponse{
			ID:             perm.ID,
			Module:         perm.Module,
			Action:         perm.Action,
			Resource:       perm.Resource,
			PermissionName: perm.PermissionName,
			Description:    perm.Description,
			IsActive:       perm.IsActive,
		})
	}

	return &response.EffectivePermissionsResponse{
		UserID:       userID,
		Permissions:  convertPermissionsToResponse(effectivePermissions),
		ModuleGroups: moduleGroups,
		Total:        len(effectivePermissions),
	}, nil
}

func (s *permissionService) BulkAssignPermissions(req request.BulkAssignPermissionsRequest) error {
	return s.permissionRepo.BulkAssignPermissions(req)
}

func (s *permissionService) CloneRolePermissions(fromRoleID, toRoleID uint, grantedBy uint) error {
	// Get source role permissions
	sourceRole, err := s.roleRepo.FindByIDWithPermissions(fromRoleID)
	if err != nil {
		return errors.New("source role not found")
	}

	// Validate target role exists
	_, err = s.roleRepo.FindByID(toRoleID)
	if err != nil {
		return errors.New("target role not found")
	}

	// Extract permission IDs
	permissionIDs := make([]uint, len(sourceRole.Permissions))
	for i, perm := range sourceRole.Permissions {
		permissionIDs[i] = perm.ID
	}

	return s.permissionRepo.AssignPermissionsToRole(toRoleID, permissionIDs, grantedBy)
}

// Helper functions
func convertPermissionsToResponse(permissions []models.Permission) []response.PermissionResponse {
	result := make([]response.PermissionResponse, len(permissions))
	for i, perm := range permissions {
		result[i] = response.PermissionResponse{
			ID:             perm.ID,
			Module:         perm.Module,
			Action:         perm.Action,
			Resource:       perm.Resource,
			PermissionName: perm.PermissionName,
			Description:    perm.Description,
			IsActive:       perm.IsActive,
		}
	}
	return result
}

func convertUserPermissionsToResponse(userPermissions []models.UserPermission) []response.UserPermissionResponse {
	result := make([]response.UserPermissionResponse, len(userPermissions))
	for i, up := range userPermissions {
		result[i] = response.UserPermissionResponse{
			ID:         up.ID,
			UserID:     up.UserID,
			Permission: response.PermissionResponse{
				ID:             up.Permission.ID,
				Module:         up.Permission.Module,
				Action:         up.Permission.Action,
				Resource:       up.Permission.Resource,
				PermissionName: up.Permission.PermissionName,
				Description:    up.Permission.Description,
				IsActive:       up.Permission.IsActive,
			},
			GrantType: up.GrantType,
			GrantedBy: up.GrantedBy,
			GrantedAt: up.GrantedAt,
			ExpiresAt: up.ExpiresAt,
			IsActive:  up.IsActive,
			Reason:    up.Reason,
		}
	}
	return result
}