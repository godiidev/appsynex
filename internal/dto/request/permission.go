// File: internal/dto/request/permission.go
// Tạo tại: internal/dto/request/permission.go
// Mục đích: Request DTOs cho Permission APIs

package request

import "time"

// Permission CRUD Requests
type CreatePermissionRequest struct {
	Module         string `json:"module" binding:"required"`
	Action         string `json:"action" binding:"required"`
	Resource       string `json:"resource"`
	PermissionName string `json:"permission_name"`
	Description    string `json:"description"`
}

type UpdatePermissionRequest struct {
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// Role Permission Requests
type AssignPermissionsToRoleRequest struct {
	RoleID        uint   `json:"role_id" binding:"required"`
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

type RemovePermissionsFromRoleRequest struct {
	RoleID        uint   `json:"role_id" binding:"required"`
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// User Permission Requests
type GrantUserPermissionRequest struct {
	UserID       uint       `json:"user_id" binding:"required"`
	PermissionID uint       `json:"permission_id" binding:"required"`
	GrantType    string     `json:"grant_type" binding:"required,oneof=GRANT DENY"` // GRANT or DENY
	GrantedBy    uint       `json:"granted_by" binding:"required"`
	ExpiresAt    *time.Time `json:"expires_at"`
	Reason       string     `json:"reason"`
}

type RevokeUserPermissionRequest struct {
	UserID       uint `json:"user_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
}

// Bulk Operations
type BulkAssignPermissionsRequest struct {
	RoleIDs       []uint `json:"role_ids" binding:"required"`
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
	GrantedBy     uint   `json:"granted_by" binding:"required"`
}

type BulkRevokePermissionsRequest struct {
	RoleIDs       []uint `json:"role_ids" binding:"required"`
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

type CloneRolePermissionsRequest struct {
	FromRoleID uint `json:"from_role_id" binding:"required"`
	ToRoleID   uint `json:"to_role_id" binding:"required"`
	GrantedBy  uint `json:"granted_by" binding:"required"`
}

// Permission Check Request
type CheckPermissionRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Module   string `json:"module" binding:"required"`
	Action   string `json:"action" binding:"required"`
	Resource string `json:"resource"`
}
