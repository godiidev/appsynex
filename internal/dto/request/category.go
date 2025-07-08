// File: internal/dto/request/category.go
// Tạo tại: internal/dto/request/category.go
// Mục đích: Định nghĩa các request DTO cho Category API

package request

type CreateCategoryRequest struct {
	CategoryName     string `json:"category_name" binding:"required"`
	ParentCategoryID *uint  `json:"parent_category_id"`
	Description      string `json:"description"`
}

type UpdateCategoryRequest struct {
	CategoryName     string `json:"category_name"`
	ParentCategoryID *uint  `json:"parent_category_id"`
	Description      string `json:"description"`
}
