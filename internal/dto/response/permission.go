// File: internal/dto/response/permission.go
// Tạo tại: internal/dto/response/permission.go
// Mục đích: Response DTOs cho Permission APIs

package response

import "time"

// Permission Responses
type PermissionResponse struct {
	ID             uint   `json:"id"`
	Module         string `json:"module"`
	Action         string `json:"action"`
	Resource       string `json:"resource"`
	PermissionName string `json:"permission_name"`
	Description    string `json:"description"`
	IsActive       bool   `json:"is_active"`
}

type PermissionsResponse struct {
	Permissions  []PermissionResponse            `json:"permissions"`
	ModuleGroups map[string][]PermissionResponse `json:"module_groups"`
	Total        int                             `json:"total"`
}

// Permission Group Responses
type PermissionGroupResponse struct {
	ID          uint                 `json:"id"`
	GroupName   string               `json:"group_name"`
	DisplayName string               `json:"display_name"`
	Description string               `json:"description"`
	Module      string               `json:"module"`
	SortOrder   int                  `json:"sort_order"`
	IsActive    bool                 `json:"is_active"`
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionGroupsResponse struct {
	Groups []PermissionGroupResponse `json:"groups"`
	Total  int                       `json:"total"`
}

// Role Permission Responses
type RolePermissionsResponse struct {
	RoleID       uint                            `json:"role_id"`
	RoleName     string                          `json:"role_name"`
	Permissions  []PermissionResponse            `json:"permissions"`
	ModuleGroups map[string][]PermissionResponse `json:"module_groups"`
	Total        int                             `json:"total"`
}

// User Permission Responses
type UserPermissionResponse struct {
	ID         uint               `json:"id"`
	UserID     uint               `json:"user_id"`
	Permission PermissionResponse `json:"permission"`
	GrantType  string             `json:"grant_type"`
	GrantedBy  uint               `json:"granted_by"`
	GrantedAt  time.Time          `json:"granted_at"`
	ExpiresAt  *time.Time         `json:"expires_at"`
	IsActive   bool               `json:"is_active"`
	Reason     string             `json:"reason"`
}

type UserPermissionsResponse struct {
	UserID            uint                     `json:"user_id"`
	Username          string                   `json:"username"`
	RolePermissions   []PermissionResponse     `json:"role_permissions"`
	DirectPermissions []UserPermissionResponse `json:"direct_permissions"`
}

// Effective Permissions Response
type EffectivePermissionsResponse struct {
	UserID       uint                            `json:"user_id"`
	Permissions  []PermissionResponse            `json:"permissions"`
	ModuleGroups map[string][]PermissionResponse `json:"module_groups"`
	Total        int                             `json:"total"`
}

// Permission Check Response
type PermissionCheckResponse struct {
	UserID    uint      `json:"user_id"`
	Module    string    `json:"module"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	HasAccess bool      `json:"has_access"`
	CheckedAt time.Time `json:"checked_at"`
}

// Bulk Operations Response
type BulkPermissionResponse struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors,omitempty"`
	Message      string   `json:"message"`
}
