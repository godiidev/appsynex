// File: internal/dto/response/category.go
// Tạo tại: internal/dto/response/category.go
// Mục đích: Định nghĩa các response DTO cho Category API

package response

import "time"

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total"`
}

type CategoryDetailResponse struct {
	ID               uint      `json:"id"`
	CategoryName     string    `json:"category_name"`
	ParentCategoryID *uint     `json:"parent_category_id"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CategoryResponse struct {
	ID               uint      `json:"id"`
	CategoryName     string    `json:"category_name"`
	ParentCategoryID *uint     `json:"parent_category_id"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
