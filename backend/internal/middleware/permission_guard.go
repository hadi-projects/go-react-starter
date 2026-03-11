package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

func PermissionGuard(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		AddToTrace(c, "PermissionGuard("+requiredPermission+")")
		permissionsInterface, exists := c.Get("permissions")
		if !exists {
			response.Error(c, http.StatusForbidden, "Permissions not found in context")
			c.Abort()
			return
		}

		permissions, ok := permissionsInterface.([]string)
		if !ok {
			response.Error(c, http.StatusForbidden, "Invalid permissions format")
			c.Abort()
			return
		}

		for _, perm := range permissions {
			if perm == requiredPermission {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Forbidden: Insufficient permissions")
		c.Abort()
	}
}
