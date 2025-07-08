// File: internal/domain/services/sample.go
// Complete and optimized version

package services

import (
	"errors"
	"math"

	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/dto/request"
	"github.com/godiidev/appsynex/internal/dto/response"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
)

type SampleService interface {
	GetSamples(req request.SampleFilterRequest) (*response.PaginatedResponse, error)
	GetSampleByID(id uint) (*response.SampleResponse, error)
	CreateSample(req request.CreateSampleRequest) (*response.SampleResponse, error)
	UpdateSample(id uint, req request.UpdateSampleRequest) (*response.SampleResponse, error)
	DeleteSample(id uint) error
}

type sampleService struct {
	sampleRepo      interfaces.SampleRepository
	productNameRepo interfaces.ProductNameRepository
	categoryRepo    interfaces.ProductCategoryRepository
}

func NewSampleService(
	sampleRepo interfaces.SampleRepository,
	productNameRepo interfaces.ProductNameRepository,
	categoryRepo interfaces.ProductCategoryRepository,
) SampleService {
	return &sampleService{
		sampleRepo:      sampleRepo,
		productNameRepo: productNameRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *sampleService) GetSamples(req request.SampleFilterRequest) (*response.PaginatedResponse, error) {
	// Set defaults for pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Build filters map
	filters := s.buildFilters(req)

	// Get samples from repository (với Preload relationships)
	samples, total, err := s.sampleRepo.FindAll(req.Page, req.Limit, req.Search, req.Category, filters)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	items := make([]interface{}, len(samples))
	for i, sample := range samples {
		items[i] = s.convertSampleToResponse(&sample)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &response.PaginatedResponse{
		Items:      items,
		TotalItems: total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *sampleService) GetSampleByID(id uint) (*response.SampleResponse, error) {
	sample, err := s.sampleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("sample not found")
	}

	return s.convertSampleToResponse(sample), nil
}

func (s *sampleService) CreateSample(req request.CreateSampleRequest) (*response.SampleResponse, error) {
	// Check if SKU already exists
	existingSample, _ := s.sampleRepo.FindBySKU(req.SKU)
	if existingSample != nil {
		return nil, errors.New("SKU already exists")
	}

	// Validate product name exists
	_, err := s.productNameRepo.FindByID(req.ProductNameID)
	if err != nil {
		return nil, errors.New("product name not found")
	}

	// Validate category exists
	_, err = s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Create sample
	sample := &models.SampleProduct{
		SKU:               req.SKU,
		ProductNameID:     req.ProductNameID,
		CategoryID:        req.CategoryID,
		Description:       req.Description,
		SampleType:        req.SampleType,
		Weight:            req.Weight,
		Width:             req.Width,
		Color:             req.Color,
		ColorCode:         req.ColorCode,
		Quality:           req.Quality,
		RemainingQuantity: req.RemainingQuantity,
		FiberContent:      req.FiberContent,
		Source:            req.Source,
		SampleLocation:    req.SampleLocation,
		Barcode:           req.Barcode,
	}

	if err := s.sampleRepo.Create(sample); err != nil {
		return nil, err
	}

	// Get complete sample with related data (Preload sẽ load relationships)
	createdSample, err := s.sampleRepo.FindByID(sample.ID)
	if err != nil {
		return nil, err
	}

	return s.convertSampleToResponse(createdSample), nil
}

func (s *sampleService) UpdateSample(id uint, req request.UpdateSampleRequest) (*response.SampleResponse, error) {
	// Get existing sample
	sample, err := s.sampleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("sample not found")
	}

	// Check if SKU is being changed and already exists
	if req.SKU != "" && req.SKU != sample.SKU {
		existingSample, _ := s.sampleRepo.FindBySKU(req.SKU)
		if existingSample != nil {
			return nil, errors.New("SKU already exists")
		}
	}

	// Update fields if provided
	if req.SKU != "" {
		sample.SKU = req.SKU
	}
	if req.ProductNameID != 0 {
		// Validate product name exists
		_, err := s.productNameRepo.FindByID(req.ProductNameID)
		if err != nil {
			return nil, errors.New("product name not found")
		}
		sample.ProductNameID = req.ProductNameID
	}
	if req.CategoryID != 0 {
		// Validate category exists
		_, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		sample.CategoryID = req.CategoryID
	}
	if req.Description != nil {
		sample.Description = *req.Description
	}
	if req.SampleType != nil {
		sample.SampleType = *req.SampleType
	}
	if req.Weight != nil {
		sample.Weight = *req.Weight
	}
	if req.Width != nil {
		sample.Width = *req.Width
	}
	if req.Color != nil {
		sample.Color = *req.Color
	}
	if req.Quality != nil {
		sample.Quality = *req.Quality
	}
	if req.RemainingQuantity != nil {
		sample.RemainingQuantity = *req.RemainingQuantity
	}
	if req.FiberContent != nil {
		sample.FiberContent = *req.FiberContent
	}
	if req.Source != nil {
		sample.Source = *req.Source
	}
	if req.SampleLocation != nil {
		sample.SampleLocation = *req.SampleLocation
	}
	if req.Barcode != nil {
		sample.Barcode = *req.Barcode
	}

	// Update sample
	if err := s.sampleRepo.Update(sample); err != nil {
		return nil, err
	}

	// Get updated sample with related data (Preload sẽ load relationships)
	updatedSample, err := s.sampleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.convertSampleToResponse(updatedSample), nil
}

func (s *sampleService) DeleteSample(id uint) error {
	// Check if sample exists
	_, err := s.sampleRepo.FindByID(id)
	if err != nil {
		return errors.New("sample not found")
	}

	return s.sampleRepo.Delete(id)
}

// Helper function to build filters map
func (s *sampleService) buildFilters(req request.SampleFilterRequest) map[string]interface{} {
	filters := make(map[string]interface{})
	if req.WeightMin > 0 {
		filters["weight_min"] = req.WeightMin
	}
	if req.WeightMax > 0 {
		filters["weight_max"] = req.WeightMax
	}
	if req.WidthMin > 0 {
		filters["width_min"] = req.WidthMin
	}
	if req.WidthMax > 0 {
		filters["width_max"] = req.WidthMax
	}
	if req.SampleType != "" {
		filters["sample_type"] = req.SampleType
	}
	if req.Color != "" {
		filters["color"] = req.Color
	}
	return filters
}

// Enhanced helper function to convert model to response DTO
func (s *sampleService) convertSampleToResponse(sample *models.SampleProduct) *response.SampleResponse {
	sampleResponse := &response.SampleResponse{
		ID:                sample.ID,
		SKU:               sample.SKU,
		ProductNameID:     sample.ProductNameID,
		CategoryID:        sample.CategoryID,
		Description:       sample.Description,
		SampleType:        sample.SampleType,
		Weight:            sample.Weight,
		Width:             sample.Width,
		Color:             sample.Color,
		ColorCode:         sample.ColorCode,
		Quality:           sample.Quality,
		RemainingQuantity: sample.RemainingQuantity,
		FiberContent:      sample.FiberContent,
		Source:            sample.Source,
		SampleLocation:    sample.SampleLocation,
		Barcode:           sample.Barcode,
		CreatedAt:         sample.CreatedAt,
		UpdatedAt:         sample.UpdatedAt,
	}

	// Auto-populate ProductName if loaded (check if relationship exists)
	if sample.ProductName.ID != 0 {
		sampleResponse.ProductName = &response.ProductNameResponse{
			ID:            sample.ProductName.ID,
			ProductNameVI: sample.ProductName.ProductNameVI,
			ProductNameEN: sample.ProductName.ProductNameEN,
			SKUParent:     sample.ProductName.SKUParent,
		}
	} else {
		// Fallback: Load product name manually if not preloaded
		if productName, err := s.productNameRepo.FindByID(sample.ProductNameID); err == nil {
			sampleResponse.ProductName = &response.ProductNameResponse{
				ID:            productName.ID,
				ProductNameVI: productName.ProductNameVI,
				ProductNameEN: productName.ProductNameEN,
				SKUParent:     productName.SKUParent,
			}
		}
	}

	// Auto-populate Category if loaded - Use existing CategoryResponse struct
	if sample.Category.ID != 0 {
		sampleResponse.Category = &response.CategoryResponse{
			ID:               sample.Category.ID,
			CategoryName:     sample.Category.CategoryName,
			ParentCategoryID: sample.Category.ParentCategoryID,
			Description:      sample.Category.Description,
			CreatedAt:        sample.Category.CreatedAt,
			UpdatedAt:        sample.Category.UpdatedAt,
		}
	} else {
		// Fallback: Load category manually if not preloaded
		if category, err := s.categoryRepo.FindByID(sample.CategoryID); err == nil {
			sampleResponse.Category = &response.CategoryResponse{
				ID:               category.ID,
				CategoryName:     category.CategoryName,
				ParentCategoryID: category.ParentCategoryID,
				Description:      category.Description,
				CreatedAt:        category.CreatedAt,
				UpdatedAt:        category.UpdatedAt,
			}
		}
	}

	return sampleResponse
}
