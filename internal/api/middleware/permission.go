// File: internal/api/middleware/permission.go
// Cập nhật tại: internal/api/middleware/permission.go
// Mục đích: Enhanced permission middleware với granular control

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"github.com/godiidev/appsynex/pkg/auth"
)

// PermissionMiddleware struct
type PermissionMiddleware struct {
	permissionRepo interfaces.PermissionRepository
}

// NewPermissionMiddleware creates new permission middleware
func NewPermissionMiddleware(permissionRepo interfaces.PermissionRepository) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionRepo: permissionRepo,
	}
}

// RequirePermission checks if user has specific permission
func (pm *PermissionMiddleware) RequirePermission(module string, action string, resource ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := userValue.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
			return
		}

		// Build resource string
		resourceStr := ""
		if len(resource) > 0 {
			resourceStr = resource[0]
		}

		// Check permission using repository
		hasPermission, err := pm.permissionRepo.UserHasPermission(claims.ID, module, action, resourceStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"module":  module,
				"action":  action,
				"resource": resourceStr,
			})
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if user has any of the specified permissions
func (pm *PermissionMiddleware) RequireAnyPermission(permissions ...PermissionCheck) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := userValue.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
			return
		}

		// Check if user has any of the permissions
		for _, perm := range permissions {
			hasPermission, err := pm.permissionRepo.UserHasPermission(claims.ID, perm.Module, perm.Action, perm.Resource)
			if err == nil && hasPermission {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	}
}

// RequireAllPermissions checks if user has all of the specified permissions
func (pm *PermissionMiddleware) RequireAllPermissions(permissions ...PermissionCheck) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := userValue.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
			return
		}

		// Check if user has all permissions
		for _, perm := range permissions {
			hasPermission, err := pm.permissionRepo.UserHasPermission(claims.ID, perm.Module, perm.Action, perm.Resource)
			if err != nil || !hasPermission {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error":   "Insufficient permissions",
					"required": perm,
				})
				return
			}
		}

		c.Next()
	}
}

// ResourceOwnerOrPermission checks if user owns the resource or has permission
func (pm *PermissionMiddleware) ResourceOwnerOrPermission(module string, action string, ownerField string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := userValue.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
			return
		}

		// Check if user is the owner of the resource
		resourceOwnerID := c.Param(ownerField)
		if resourceOwnerID != "" && resourceOwnerID == string(rune(claims.ID)) {
			c.Next()
			return
		}

		// Check permission
		hasPermission, err := pm.permissionRepo.UserHasPermission(claims.ID, module, action, "")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}

// RoleBasedAccess checks if user has specific role
func RoleBasedAccess(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := userValue.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
			return
		}

		// Check if user has any of the allowed roles
		for _, userRole := range claims.Roles {
			for _, allowedRole := range allowedRoles {
				if strings.EqualFold(userRole, allowedRole) {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied - insufficient role"})
	}
}

// SuperAdminOnly restricts access to super admin only
func SuperAdminOnly() gin.HandlerFunc {
	return RoleBasedAccess("SUPER_ADMIN", "ADMIN")
}

// PermissionCheck struct for complex permission checking
type PermissionCheck struct {
	Module   string
	Action   string
	Resource string
}

// Helper functions for common permission patterns

// CanView checks VIEW permission for a module
func CanView(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		pm := getPermissionMiddleware(c)
		pm.RequirePermission(module, "VIEW")(c)
	}
}

// CanCreate checks CREATE permission for a module
func CanCreate(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		pm := getPermissionMiddleware(c)
		pm.RequirePermission(module, "CREATE")(c)
	}
}

// CanUpdate checks UPDATE permission for a module
func CanUpdate(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		pm := getPermissionMiddleware(c)
		pm.RequirePermission(module, "UPDATE")(c)
	}
}

// CanDelete checks DELETE permission for a module
func CanDelete(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		pm := getPermissionMiddleware(c)
		pm.RequirePermission(module, "DELETE")(c)
	}
}

// CanManage checks if user can perform any CRUD operation on a module
func CanManage(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		pm := getPermissionMiddleware(c)
		pm.RequireAnyPermission(
			PermissionCheck{Module: module, Action: "VIEW"},
			PermissionCheck{Module: module, Action: "CREATE"},
			PermissionCheck{Module: module, Action: "UPDATE"},
			PermissionCheck{Module: module, Action: "DELETE"},
		)(c)
	}
}

// Helper function to get permission middleware from context
func getPermissionMiddleware(c *gin.Context) *PermissionMiddleware {
	// This should be injected into context during setup
	pm, exists := c.Get("permissionMiddleware")
	if !exists {
		// Fallback - this should not happen in production
		panic("Permission middleware not found in context")
	}
	return pm.(*PermissionMiddleware)
}

// Middleware to inject permission middleware into context
func InjectPermissionMiddleware(permissionRepo interfaces.PermissionRepository) gin.HandlerFunc {
	pm := NewPermissionMiddleware(permissionRepo)
	return func(c *gin.Context) {
		c.Set("permissionMiddleware", pm)
		c.Next()
	}
}

// Legacy HasPermission function for backward compatibility
func HasPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := userValue.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
			return
		}

		// Check permission using legacy method
		hasPermission := false
		for _, perm := range claims.Permissions {
			if perm.Name == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}