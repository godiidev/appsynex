package interfaces

import "github.com/godiidev/appsynex/internal/domain/models"

type ProductNameRepository interface {
	FindAll() ([]models.ProductName, error)
	FindByID(id uint) (*models.ProductName, error)
	Create(productName *models.ProductName) error
	Update(productName *models.ProductName) error
	Delete(id uint) error
}
