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

---

// File: internal/repository/mysql/permission.go
// Tạo tại: internal/repository/mysql/permission.go
// Mục đích: MySQL implementation của Permission Repository

package mysql

import (
	"time"

	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/dto/request"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) interfaces.PermissionRepository {
	return &permissionRepository{db: db}
}

// Permission CRUD
func (r *permissionRepository) FindAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("is_active = ?", true).Order("module ASC, action ASC").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("permission_name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByModule(module string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("module = ? AND is_active = ?", module, true).Order("action ASC").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

func (r *permissionRepository) IsPermissionInUse(id uint) (bool, error) {
	var count int64
	
	// Check in role_permissions
	err := r.db.Model(&models.RolePermission{}).Where("permission_id = ? AND is_active = ?", id, true).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	
	// Check in user_permissions
	err = r.db.Model(&models.UserPermission{}).Where("permission_id = ? AND is_active = ?", id, true).Count(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// Permission Groups
func (r *permissionRepository) FindAllGroups() ([]models.PermissionGroup, error) {
	var groups []models.PermissionGroup
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&groups).Error
	return groups, err
}

func (r *permissionRepository) CreateGroup(group *models.PermissionGroup) error {
	return r.db.Create(group).Error
}

func (r *permissionRepository) UpdateGroup(group *models.PermissionGroup) error {
	return r.db.Save(group).Error
}

func (r *permissionRepository) DeleteGroup(id uint) error {
	return r.db.Delete(&models.PermissionGroup{}, id).Error
}

// Role Permissions
func (r *permissionRepository) AssignPermissionsToRole(roleID uint, permissionIDs []uint, grantedBy uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Remove existing permissions
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new permissions
	for _, permissionID := range permissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
			GrantedBy:    grantedBy,
			GrantedAt:    time.Now(),
			IsActive:     true,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *permissionRepository) RemovePermissionsFromRole(roleID uint, permissionIDs []uint) error {
	return r.db.Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).Delete(&models.RolePermission{}).Error
}

func (r *permissionRepository) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND role_permissions.is_active = ? AND permissions.is_active = ?", roleID, true, true).
		Find(&permissions).Error
	return permissions, err
}

// User Permissions (Direct)
func (r *permissionRepository) GrantUserPermission(req request.GrantUserPermissionRequest) error {
	// Check if permission already exists
	var existing models.UserPermission
	err := r.db.Where("user_id = ? AND permission_id = ?", req.UserID, req.PermissionID).First(&existing).Error
	
	if err == nil {
		// Update existing
		existing.GrantType = req.GrantType
		existing.GrantedBy = req.GrantedBy
		existing.GrantedAt = time.Now()
		existing.ExpiresAt = req.ExpiresAt
		existing.IsActive = true
		existing.Reason = req.Reason
		return r.db.Save(&existing).Error
	}

	// Create new
	userPermission := models.UserPermission{
		UserID:       req.UserID,
		PermissionID: req.PermissionID,
		GrantType:    req.GrantType,
		GrantedBy:    req.GrantedBy,
		GrantedAt:    time.Now(),
		ExpiresAt:    req.ExpiresAt,
		IsActive:     true,
		Reason:       req.Reason,
	}
	return r.db.Create(&userPermission).Error
}

func (r *permissionRepository) RevokeUserPermission(userID uint, permissionID uint) error {
	return r.db.Where("user_id = ? AND permission_id = ?", userID, permissionID).Delete(&models.UserPermission{}).Error
}

func (r *permissionRepository) GetUserDirectPermissions(userID uint) ([]models.UserPermission, error) {
	var userPermissions []models.UserPermission
	err := r.db.Preload("Permission").
		Where("user_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)", userID, true, time.Now()).
		Find(&userPermissions).Error
	return userPermissions, err
}

func (r *permissionRepository) GetUserRolePermissions(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND role_permissions.is_active = ? AND permissions.is_active = ?", userID, true, true).
		Where("role_permissions.expires_at IS NULL OR role_permissions.expires_at > ?", time.Now()).
		Distinct().
		Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) GetUserEffectivePermissions(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	
	// Union query to get both role-based and direct permissions
	query := `
		SELECT DISTINCT p.* FROM permissions p
		WHERE p.is_active = true AND p.id IN (
			-- Role-based permissions
			SELECT rp.permission_id 
			FROM role_permissions rp
			INNER JOIN user_roles ur ON rp.role_id = ur.role_id
			WHERE ur.user_id = ? 
			AND rp.is_active = true
			AND (rp.expires_at IS NULL OR rp.expires_at > ?)
			
			UNION
			
			-- Direct user permissions (GRANT only, not DENY)
			SELECT up.permission_id
			FROM user_permissions up
			WHERE up.user_id = ?
			AND up.is_active = true
			AND up.grant_type = 'GRANT'
			AND (up.expires_at IS NULL OR up.expires_at > ?)
		)
		-- Exclude denied permissions
		AND p.id NOT IN (
			SELECT up2.permission_id
			FROM user_permissions up2
			WHERE up2.user_id = ?
			AND up2.is_active = true
			AND up2.grant_type = 'DENY'
			AND (up2.expires_at IS NULL OR up2.expires_at > ?)
		)
		ORDER BY p.module ASC, p.action ASC
	`
	
	now := time.Now()
	err := r.db.Raw(query, userID, now, userID, now, userID, now).Find(&permissions).Error
	return permissions, err
}

// Permission Checking
func (r *permissionRepository) UserHasPermission(userID uint, module string, action string, resource string) (bool, error) {
	var count int64
	
	// Build permission name
	permissionName := module + "_" + action
	if resource != "" {
		permissionName += "_" + resource
	}
	
	// Check effective permissions
	query := `
		SELECT COUNT(DISTINCT p.id) FROM permissions p
		WHERE p.is_active = true 
		AND (p.permission_name = ? OR (p.module = ? AND p.action = ?))
		AND p.id IN (
			-- Role-based permissions
			SELECT rp.permission_id 
			FROM role_permissions rp
			INNER JOIN user_roles ur ON rp.role_id = ur.role_id
			WHERE ur.user_id = ? 
			AND rp.is_active = true
			AND (rp.expires_at IS NULL OR rp.expires_at > ?)
			
			UNION
			
			-- Direct user permissions (GRANT only)
			SELECT up.permission_id
			FROM user_permissions up
			WHERE up.user_id = ?
			AND up.is_active = true
			AND up.grant_type = 'GRANT'
			AND (up.expires_at IS NULL OR up.expires_at > ?)
		)
		-- Exclude denied permissions
		AND p.id NOT IN (
			SELECT up2.permission_id
			FROM user_permissions up2
			WHERE up2.user_id = ?
			AND up2.is_active = true
			AND up2.grant_type = 'DENY'
			AND (up2.expires_at IS NULL OR up2.expires_at > ?)
		)
	`
	
	now := time.Now()
	err := r.db.Raw(query, permissionName, module, action, userID, now, userID, now, userID, now).Count(&count).Error
	
	return count > 0, err
}

// Bulk Operations
func (r *permissionRepository) BulkAssignPermissions(req request.BulkAssignPermissionsRequest) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, roleID := range req.RoleIDs {
		if err := r.AssignPermissionsToRole(roleID, req.PermissionIDs, req.GrantedBy); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *permissionRepository) BulkRevokePermissions(req request.BulkRevokePermissionsRequest) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, roleID := range req.RoleIDs {
		if err := r.RemovePermissionsFromRole(roleID, req.PermissionIDs); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}