package interfaces

import (
	"github.com/godiidev/appsynex/internal/domain/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]models.Role, error)
	FindByID(id uint) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint) error
	FindByIDWithPermissions(id uint) (*models.Role, error)
	AssignPermissions(roleID uint, permissions []models.RolePermission) error
	RemoveAllPermissions(roleID uint) error
	GetDB() *gorm.DB
}
