package interfaces

import "github.com/godiidev/appsynex/internal/domain/models"

type UserRepository interface {
	FindAll(page, limit int, search string) ([]models.User, int64, error)
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	FindByIDWithRoles(id uint) (*models.User, error)
	AssignRoles(userID uint, roleIDs []uint) error
}
