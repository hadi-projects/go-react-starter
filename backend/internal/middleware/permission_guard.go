package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type PermissionGuard struct {
	cache cache.CacheService
	repo  repository.PermissionRepository
}

func NewPermissionGuard(cache cache.CacheService, repo repository.PermissionRepository) *PermissionGuard {
	return &PermissionGuard{
		cache: cache,
		repo:  repo,
	}
}

func (g *PermissionGuard) Check(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		AddToTrace(c, "PermissionGuard("+requiredPermission+")")

		role, exists := c.Get("role")
		if exists && role != nil {
			if roleStr, ok := role.(string); ok && roleStr == "admin" {
				c.Next()
				return
			}
		}

		permissionsMaskInterface, exists := c.Get("permissions_mask")
		if !exists {
			response.Error(c, http.StatusForbidden, "Permissions mask not found in context")
			c.Abort()
			return
		}

		userMask, ok := permissionsMaskInterface.(uint64)
		if !ok {
			response.Error(c, http.StatusForbidden, "Invalid permissions mask format")
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		
		cacheKey := "perm_id:" + requiredPermission
		var permID uint

		// Try to get from cache
		if err := g.cache.Get(ctx, cacheKey, &permID); err != nil {
			// Not in cache, get from db (FindByPermissionName logic depends on repo, wait! FindByName returns a single Permission)
			perm, err := g.repo.FindByName(ctx, requiredPermission)
			if err != nil {
				logger.SystemLogger.Error().Err(err).Msg("Failed to find permission by name in guard")
				response.Error(c, http.StatusInternalServerError, "Failed to verify permission")
				c.Abort()
				return
			}
			permID = perm.ID
			// Set cache for 24 hours
			g.cache.Set(ctx, cacheKey, permID, 24*time.Hour)
		}

		if permID == 0 || permID > 64 {
			response.Error(c, http.StatusForbidden, "Invalid permission configuration")
			c.Abort()
			return
		}

		requiredMask := uint64(1) << (permID - 1)

		if (userMask & requiredMask) != 0 {
			c.Next()
			return
		}

		response.Error(c, http.StatusForbidden, "Forbidden: Insufficient permissions")
		c.Abort()
	}
}
