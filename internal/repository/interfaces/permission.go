// File: internal/repository/interfaces/permission.go
// Tạo tại: internal/repository/interfaces/permission.go
// Mục đích: Interface cho Permission Repository

package interfaces

import (
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/dto/request"
)

type PermissionRepository interface {
	// Permission CRUD
	FindAll() ([]models.Permission, error)
	FindByID(id uint) (*models.Permission, error)
	FindByName(name string) (*models.Permission, error)
	FindByModule(module string) ([]models.Permission, error)
	Create(permission *models.Permission) error
	Update(permission *models.Permission) error
	Delete(id uint) error
	IsPermissionInUse(id uint) (bool, error)

	// Permission Groups
	FindAllGroups() ([]models.PermissionGroup, error)
	CreateGroup(group *models.PermissionGroup) error
	UpdateGroup(group *models.PermissionGroup) error
	DeleteGroup(id uint) error

	// Role Permissions
	AssignPermissionsToRole(roleID uint, permissionIDs []uint, grantedBy uint) error
	RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]models.Permission, error)

	// User Permissions (Direct)
	GrantUserPermission(req request.GrantUserPermissionRequest) error
	RevokeUserPermission(userID uint, permissionID uint) error
	GetUserDirectPermissions(userID uint) ([]models.UserPermission, error)
	GetUserRolePermissions(userID uint) ([]models.Permission, error)
	GetUserEffectivePermissions(userID uint) ([]models.Permission, error)

	// Permission Checking
	UserHasPermission(userID uint, module string, action string, resource string) (bool, error)

	// Bulk Operations
	BulkAssignPermissions(req request.BulkAssignPermissionsRequest) error
	BulkRevokePermissions(req request.BulkRevokePermissionsRequest) error
}
