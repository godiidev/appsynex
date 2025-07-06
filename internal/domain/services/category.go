// File: internal/domain/services/category.go
// Tạo tại: internal/domain/services/category.go
// Mục đích: Service xử lý logic nghiệp vụ cho Product Category (CRUD, validate parent-child relationships)

package services

import (
	"errors"

	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/dto/request"
	"github.com/godiidev/appsynex/internal/dto/response"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
)

type CategoryService interface {
	GetAllCategories() (*response.CategoriesResponse, error)
	GetCategoryByID(id uint) (*response.CategoryDetailResponse, error)
	CreateCategory(req request.CreateCategoryRequest) (*response.CategoryDetailResponse, error)
	UpdateCategory(id uint, req request.UpdateCategoryRequest) (*response.CategoryDetailResponse, error)
	DeleteCategory(id uint) error
}

type categoryService struct {
	categoryRepo interfaces.ProductCategoryRepository
}

func NewCategoryService(categoryRepo interfaces.ProductCategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetAllCategories() (*response.CategoriesResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	items := make([]response.CategoryResponse, len(categories))
	for i, category := range categories {
		items[i] = response.CategoryResponse{
			ID:               category.ID,
			CategoryName:     category.CategoryName,
			ParentCategoryID: category.ParentCategoryID,
			Description:      category.Description,
			CreatedAt:        category.CreatedAt,
			UpdatedAt:        category.UpdatedAt,
		}
	}

	return &response.CategoriesResponse{
		Categories: items,
		Total:      len(items),
	}, nil
}

func (s *categoryService) GetCategoryByID(id uint) (*response.CategoryDetailResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	return convertCategoryToDetailResponse(category), nil
}

func (s *categoryService) CreateCategory(req request.CreateCategoryRequest) (*response.CategoryDetailResponse, error) {
	// Validate parent category if provided
	if req.ParentCategoryID != nil {
		_, err := s.categoryRepo.FindByID(*req.ParentCategoryID)
		if err != nil {
			return nil, errors.New("parent category not found")
		}
	}

	// Create category
	category := &models.ProductCategory{
		CategoryName:     req.CategoryName,
		ParentCategoryID: req.ParentCategoryID,
		Description:      req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	// Get complete category
	createdCategory, err := s.categoryRepo.FindByID(category.ID)
	if err != nil {
		return nil, err
	}

	return convertCategoryToDetailResponse(createdCategory), nil
}

func (s *categoryService) UpdateCategory(id uint, req request.UpdateCategoryRequest) (*response.CategoryDetailResponse, error) {
	// Get existing category
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Validate parent category if provided
	if req.ParentCategoryID != nil {
		// Check for circular reference
		if *req.ParentCategoryID == id {
			return nil, errors.New("category cannot be its own parent")
		}

		_, err := s.categoryRepo.FindByID(*req.ParentCategoryID)
		if err != nil {
			return nil, errors.New("parent category not found")
		}
	}

	// Update fields if provided
	if req.CategoryName != "" {
		category.CategoryName = req.CategoryName
	}
	if req.ParentCategoryID != nil {
		category.ParentCategoryID = req.ParentCategoryID
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	// Update category
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	// Get updated category
	updatedCategory, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return convertCategoryToDetailResponse(updatedCategory), nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	// Check if category exists
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	// TODO: Check if category has products or samples before deleting
	// This would require additional repository methods

	return s.categoryRepo.Delete(id)
}

// Helper function to convert model to response DTO
func convertCategoryToDetailResponse(category *models.ProductCategory) *response.CategoryDetailResponse {
	return &response.CategoryDetailResponse{
		ID:               category.ID,
		CategoryName:     category.CategoryName,
		ParentCategoryID: category.ParentCategoryID,
		Description:      category.Description,
		CreatedAt:        category.CreatedAt,
		UpdatedAt:        category.UpdatedAt,
	}
}