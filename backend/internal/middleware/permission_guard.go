package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PermissionGuard(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionsInterface, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permissions not found in context"})
			c.Abort()
			return
		}

		permissions, ok := permissionsInterface.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
			c.Abort()
			return
		}

		for _, perm := range permissions {
			if perm == requiredPermission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
		c.Abort()
	}
}
