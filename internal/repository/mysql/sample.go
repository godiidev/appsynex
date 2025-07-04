package mysql

import (
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"gorm.io/gorm"
)

type sampleRepository struct {
	db *gorm.DB
}

func NewSampleRepository(db *gorm.DB) interfaces.SampleRepository {
	return &sampleRepository{db: db}
}

func (r *sampleRepository) FindAll(page, limit int, search, category string, filters map[string]interface{}) ([]models.SampleProduct, int64, error) {
	var samples []models.SampleProduct
	var count int64

	query := r.db.Model(&models.SampleProduct{})

	// Apply search
	if search != "" {
		query = query.Joins("JOIN product_names ON sample_products.product_name_id = product_names.id").
			Where("sample_products.sku LIKE ? OR product_names.product_name_vi LIKE ? OR product_names.product_name_en LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Apply category filter
	if category != "" {
		query = query.Where("category_id = ?", category)
	}

	// Apply other filters
	for key, value := range filters {
		switch key {
		case "weight_min":
			query = query.Where("weight >= ?", value)
		case "weight_max":
			query = query.Where("weight <= ?", value)
		case "width_min":
			query = query.Where("width >= ?", value)
		case "width_max":
			query = query.Where("width <= ?", value)
		case "sample_type":
			query = query.Where("sample_type = ?", value)
		case "color":
			query = query.Where("color LIKE ?", "%"+value.(string)+"%")
		}
	}

	// Count total
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Preload("ProductName").
		Preload("Category").
		Find(&samples).Error

	return samples, count, err
}

func (r *sampleRepository) FindByID(id uint) (*models.SampleProduct, error) {
	var sample models.SampleProduct
	if err := r.db.Preload("ProductName").Preload("Category").First(&sample, id).Error; err != nil {
		return nil, err
	}
	return &sample, nil
}

func (r *sampleRepository) FindBySKU(sku string) (*models.SampleProduct, error) {
	var sample models.SampleProduct
	if err := r.db.Where("sku = ?", sku).First(&sample).Error; err != nil {
		return nil, err
	}
	return &sample, nil
}

func (r *sampleRepository) Create(sample *models.SampleProduct) error {
	return r.db.Create(sample).Error
}

func (r *sampleRepository) Update(sample *models.SampleProduct) error {
	return r.db.Save(sample).Error
}

func (r *sampleRepository) Delete(id uint) error {
	return r.db.Delete(&models.SampleProduct{}, id).Error
}
