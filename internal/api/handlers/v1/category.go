// File: internal/api/handlers/v1/category.go
// Tạo tại: internal/api/handlers/v1/category.go  
// Mục đích: Handler xử lý các API liên quan đến quản lý danh mục sản phẩm (CRUD product_categories)

package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/godiidev/appsynex/internal/domain/services"
	"github.com/godiidev/appsynex/internal/dto/request"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetAll godoc
// @Summary     Get all product categories
// @Description Get a list of product categories
// @Tags        categories
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} response.CategoriesResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetByID godoc
// @Summary     Get category by ID
// @Description Get a product category by ID
// @Tags        categories
// @Accept      json
// @Produce     json
// @Param       id path int true "Category ID"
// @Security    BearerAuth
// @Success     200 {object} response.CategoryDetailResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Create godoc
// @Summary     Create a new category
// @Description Create a new product category
// @Tags        categories
// @Accept      json
// @Produce     json
// @Param       category body request.CreateCategoryRequest true "Category to create"
// @Security    BearerAuth
// @Success     201 {object} response.CategoryDetailResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Update godoc
// @Summary     Update a category
// @Description Update a product category by ID
// @Tags        categories
// @Accept      json
// @Produce     json
// @Param       id path int true "Category ID"
// @Param       category body request.UpdateCategoryRequest true "Category data to update"
// @Security    BearerAuth
// @Success     200 {object} response.CategoryDetailResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete godoc
// @Summary     Delete a category
// @Description Delete a product category by ID
// @Tags        categories
// @Accept      json
// @Produce     json
// @Param       id path int true "Category ID"
// @Security    BearerAuth
// @Success     204 {object} nil
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}