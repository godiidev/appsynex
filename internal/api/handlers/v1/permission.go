// File: internal/api/handlers/v1/permission.go
// Tạo tại: internal/api/handlers/v1/permission.go
// Mục đích: Handler xử lý các API liên quan đến Permission management

package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/godiidev/appsynex/internal/domain/services"
	"github.com/godiidev/appsynex/internal/dto/request"
	"github.com/godiidev/appsynex/pkg/auth"
)

type PermissionHandler struct {
	permissionService services.PermissionService
}

func NewPermissionHandler(permissionService services.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// GetAllPermissions godoc
// @Summary     Get all permissions
// @Description Get all permissions grouped by module
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} response.PermissionsResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions [get]
func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.permissionService.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetPermissionsByModule godoc
// @Summary     Get permissions by module
// @Description Get permissions filtered by module
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       module path string true "Module name"
// @Security    BearerAuth
// @Success     200 {object} response.PermissionsResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions/module/{module} [get]
func (h *PermissionHandler) GetPermissionsByModule(c *gin.Context) {
	module := c.Param("module")
	if module == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Module is required"})
		return
	}

	permissions, err := h.permissionService.GetPermissionsByModule(module)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetPermissionGroups godoc
// @Summary     Get permission groups
// @Description Get all permission groups with their permissions
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} response.PermissionGroupsResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions/groups [get]
func (h *PermissionHandler) GetPermissionGroups(c *gin.Context) {
	groups, err := h.permissionService.GetPermissionGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// CreatePermission godoc
// @Summary     Create new permission
// @Description Create a new system permission
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       permission body request.CreatePermissionRequest true "Permission to create"
// @Security    BearerAuth
// @Success     201 {object} response.PermissionResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req request.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := h.permissionService.CreatePermission(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

// UpdatePermission godoc
// @Summary     Update permission
// @Description Update a permission by ID
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       id path int true "Permission ID"
// @Param       permission body request.UpdatePermissionRequest true "Permission data to update"
// @Security    BearerAuth
// @Success     200 {object} response.PermissionResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req request.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := h.permissionService.UpdatePermission(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// DeletePermission godoc
// @Summary     Delete permission
// @Description Delete a permission by ID
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       id path int true "Permission ID"
// @Security    BearerAuth
// @Success     204 {object} nil
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.permissionService.DeletePermission(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignPermissionsToRole godoc
// @Summary     Assign permissions to role
// @Description Assign multiple permissions to a role
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       roleId path int true "Role ID"
// @Param       permissions body request.AssignPermissionsToRoleRequest true "Permissions to assign"
// @Security    BearerAuth
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /roles/{roleId}/permissions [post]
func (h *PermissionHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	var req request.AssignPermissionsToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user from context
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	claims := userClaims.(*auth.JWTClaims)

	req.RoleID = uint(roleID)
	if err := h.permissionService.AssignPermissionsToRole(req.RoleID, req.PermissionIDs, claims.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions assigned successfully"})
}

// RemovePermissionsFromRole godoc
// @Summary     Remove permissions from role
// @Description Remove multiple permissions from a role
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       roleId path int true "Role ID"
// @Param       permissions body request.RemovePermissionsFromRoleRequest true "Permissions to remove"
// @Security    BearerAuth
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /roles/{roleId}/permissions [delete]
func (h *PermissionHandler) RemovePermissionsFromRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	var req request.RemovePermissionsFromRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.RoleID = uint(roleID)
	if err := h.permissionService.RemovePermissionsFromRole(req.RoleID, req.PermissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions removed successfully"})
}

// GetRolePermissions godoc
// @Summary     Get role permissions
// @Description Get all permissions assigned to a role
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       roleId path int true "Role ID"
// @Security    BearerAuth
// @Success     200 {object} response.RolePermissionsResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /roles/{roleId}/permissions [get]
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	permissions, err := h.permissionService.GetRolePermissions(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GrantUserPermission godoc
// @Summary     Grant direct permission to user
// @Description Grant a specific permission directly to a user
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       userId path int true "User ID"
// @Param       permission body request.GrantUserPermissionRequest true "Permission to grant"
// @Security    BearerAuth
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{userId}/permissions [post]
func (h *PermissionHandler) GrantUserPermission(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req request.GrantUserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user from context
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	claims := userClaims.(*auth.JWTClaims)

	req.UserID = uint(userID)
	req.GrantedBy = claims.ID

	if err := h.permissionService.GrantUserPermission(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission granted successfully"})
}

// RevokeUserPermission godoc
// @Summary     Revoke direct permission from user
// @Description Revoke a specific permission directly from a user
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       userId path int true "User ID"
// @Param       permissionId path int true "Permission ID"
// @Security    BearerAuth
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{userId}/permissions/{permissionId} [delete]
func (h *PermissionHandler) RevokeUserPermission(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	permissionID, err := strconv.ParseUint(c.Param("permissionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID format"})
		return
	}

	if err := h.permissionService.RevokeUserPermission(uint(userID), uint(permissionID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked successfully"})
}

// GetUserPermissions godoc
// @Summary     Get user permissions
// @Description Get all permissions for a user (both role-based and direct)
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       userId path int true "User ID"
// @Security    BearerAuth
// @Success     200 {object} response.UserPermissionsResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{userId}/permissions [get]
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	permissions, err := h.permissionService.GetUserPermissions(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetUserEffectivePermissions godoc
// @Summary     Get user effective permissions
// @Description Get all effective permissions for a user (merged role-based and direct permissions)
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       userId path int true "User ID"
// @Security    BearerAuth
// @Success     200 {object} response.EffectivePermissionsResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{userId}/effective-permissions [get]
func (h *PermissionHandler) GetUserEffectivePermissions(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	permissions, err := h.permissionService.GetUserEffectivePermissions(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// CheckUserPermission godoc
// @Summary     Check user permission
// @Description Check if a user has a specific permission
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       permission body request.CheckPermissionRequest true "Permission to check"
// @Security    BearerAuth
// @Success     200 {object} response.PermissionCheckResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions/check [post]
func (h *PermissionHandler) CheckUserPermission(c *gin.Context) {
	var req request.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasPermission, err := h.permissionService.UserHasPermission(req.UserID, req.Module, req.Action, req.Resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"user_id":    req.UserID,
		"module":     req.Module,
		"action":     req.Action,
		"resource":   req.Resource,
		"has_access": hasPermission,
		"checked_at": "time.Now()",
	}

	c.JSON(http.StatusOK, response)
}

// BulkAssignPermissions godoc
// @Summary     Bulk assign permissions
// @Description Assign multiple permissions to multiple roles
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       bulk body request.BulkAssignPermissionsRequest true "Bulk assignment data"
// @Security    BearerAuth
// @Success     200 {object} response.BulkPermissionResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions/bulk-assign [post]
func (h *PermissionHandler) BulkAssignPermissions(c *gin.Context) {
	var req request.BulkAssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user from context
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	claims := userClaims.(*auth.JWTClaims)
	req.GrantedBy = claims.ID

	if err := h.permissionService.BulkAssignPermissions(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success_count": len(req.RoleIDs),
		"message":       "Permissions assigned successfully to all roles",
	})
}

// CloneRolePermissions godoc
// @Summary     Clone role permissions
// @Description Copy all permissions from one role to another
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Param       clone body request.CloneRolePermissionsRequest true "Clone data"
// @Security    BearerAuth
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /permissions/clone-role [post]
func (h *PermissionHandler) CloneRolePermissions(c *gin.Context) {
	var req request.CloneRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user from context
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	claims := userClaims.(*auth.JWTClaims)
	req.GrantedBy = claims.ID

	if err := h.permissionService.CloneRolePermissions(req.FromRoleID, req.ToRoleID, req.GrantedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role permissions cloned successfully"})
}
