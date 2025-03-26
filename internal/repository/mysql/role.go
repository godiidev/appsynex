package mysql

import (
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) interfaces.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *roleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *roleRepository) FindByIDWithPermissions(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) AssignPermissions(roleID uint, permissions []models.RolePermission) error {
	// Begin transaction
	tx := r.db.Begin()

	// Remove existing permissions
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new permissions
	for _, permission := range permissions {
		permission.RoleID = roleID
		if err := tx.Create(&permission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

func (r *roleRepository) RemoveAllPermissions(roleID uint) error {
	return r.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}
