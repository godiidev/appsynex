package interfaces

import "github.com/godiidev/appsynex/internal/domain/models"

type SampleRepository interface {
	FindAll(page, limit int, search, category string, filters map[string]interface{}) ([]models.SampleProduct, int64, error)
	FindByID(id uint) (*models.SampleProduct, error)
	FindBySKU(sku string) (*models.SampleProduct, error)
	Create(sample *models.SampleProduct) error
	Update(sample *models.SampleProduct) error
	Delete(id uint) error
}
