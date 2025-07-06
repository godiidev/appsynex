// File: internal/dto/request/user.go
// Tạo tại: internal/dto/request/user.go
// Mục đích: Định nghĩa các request DTO cho User API

package request

type UserFilterRequest struct {
	Page   int    `form:"page" json:"page"`
	Limit  int    `form:"limit" json:"limit"`
	Search string `form:"search" json:"search"`
}

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required,min=6"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    *string `json:"phone"`
	RoleIDs  []uint  `json:"role_ids"`
}

type UpdateUserRequest struct {
	Username      string  `json:"username"`
	Password      string  `json:"password,omitempty"`
	Email         string  `json:"email"`
	Phone         *string `json:"phone"`
	AccountStatus string  `json:"account_status"`
	RoleIDs       []uint  `json:"role_ids"`
}

type AssignRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// File: internal/dto/request/category.go
// Tạo tại: internal/dto/request/category.go
// Mục đích: Định nghĩa các request DTO cho Category API

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