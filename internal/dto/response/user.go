// File: internal/dto/response/user.go
// Tạo tại: internal/dto/response/user.go
// Mục đích: Định nghĩa các response DTO cho User API

package response

import "time"

type UserDetailResponse struct {
	ID            uint           `json:"id"`
	Username      string         `json:"username"`
	Email         string         `json:"email"`
	Phone         *string        `json:"phone"`
	LastLogin     *time.Time     `json:"last_login"`
	AccountStatus string         `json:"account_status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Roles         []RoleResponse `json:"roles"`
}

type RoleResponse struct {
	ID          uint   `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}
