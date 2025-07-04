package mysql

import (
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"gorm.io/gorm"
)

type productCategoryRepository struct {
	db *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) interfaces.ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func (r *productCategoryRepository) FindAll() ([]models.ProductCategory, error) {
	var categories []models.ProductCategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *productCategoryRepository) FindByID(id uint) (*models.ProductCategory, error) {
	var category models.ProductCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *productCategoryRepository) Create(category *models.ProductCategory) error {
	return r.db.Create(category).Error
}

func (r *productCategoryRepository) Update(category *models.ProductCategory) error {
	return r.db.Save(category).Error
}

func (r *productCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductCategory{}, id).Error
}
