package mysql

import (
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"gorm.io/gorm"
)

type productNameRepository struct {
	db *gorm.DB
}

func NewProductNameRepository(db *gorm.DB) interfaces.ProductNameRepository {
	return &productNameRepository{db: db}
}

func (r *productNameRepository) FindAll() ([]models.ProductName, error) {
	var productNames []models.ProductName
	if err := r.db.Find(&productNames).Error; err != nil {
		return nil, err
	}
	return productNames, nil
}

func (r *productNameRepository) FindByID(id uint) (*models.ProductName, error) {
	var productName models.ProductName
	if err := r.db.First(&productName, id).Error; err != nil {
		return nil, err
	}
	return &productName, nil
}

func (r *productNameRepository) Create(productName *models.ProductName) error {
	return r.db.Create(productName).Error
}

func (r *productNameRepository) Update(productName *models.ProductName) error {
	return r.db.Save(productName).Error
}

func (r *productNameRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductName{}, id).Error
}
