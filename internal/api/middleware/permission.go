package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godiidev/appsynex/pkg/auth"
)

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

		// Check permission
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
