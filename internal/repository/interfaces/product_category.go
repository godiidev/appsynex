package interfaces

import "github.com/godiidev/appsynex/internal/domain/models"

type ProductCategoryRepository interface {
	FindAll() ([]models.ProductCategory, error)
	FindByID(id uint) (*models.ProductCategory, error)
	Create(category *models.ProductCategory) error
	Update(category *models.ProductCategory) error
	Delete(id uint) error
}
