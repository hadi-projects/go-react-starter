package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

func APIKeyMiddleware(svc service.ApiKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.Next() // Continue to check other auth methods (JWT)
			return
		}

		key, mask, err := svc.Validate(c.Request.Context(), apiKey)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired API Key")
			c.Abort()
			return
		}

		// IP Whitelisting Check
		if key.AllowedIPs != "" {
			clientIP := c.ClientIP()
			allowed := false
			ips := strings.Split(key.AllowedIPs, ",")
			for _, ip := range ips {
				if strings.TrimSpace(ip) == clientIP {
					allowed = true
					break
				}
			}
			if !allowed {
				logger.SystemLogger.Warn().
					Str("client_ip", clientIP).
					Uint("api_key_id", key.ID).
					Msg("API Key used from unauthorized IP")
				response.Error(c, http.StatusForbidden, "IP not authorized for this API Key")
				c.Abort()
				return
			}
		}

		// Inject into context to simulate a logged-in user
		c.Set("user_id", key.UserID)
		c.Set("role", key.Role.Name)
		c.Set("permissions_mask", mask)
		c.Set("api_key_id", key.ID)
		c.Set("is_api_auth", true)

		// Update request context for logger (including ApiKeyID)
		ctx := context.WithValue(c.Request.Context(), logger.CtxKeyUserID, key.UserID)
		ctx = context.WithValue(ctx, logger.CtxKeyApiKeyID, key.ID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
