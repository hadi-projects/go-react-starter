package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		AddToTrace(c, "AuthMiddleware")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			logger.SystemLogger.Info().Msgf("DEBUG: AuthMiddleware claims=%v", claims)
			var userID uint
			var userEmail string

			// Cast sub to uint
			if subFloat, ok := claims["sub"].(float64); ok {
				userID = uint(subFloat)
				c.Set("user_id", userID)
			}
			if email, ok := claims["email"].(string); ok {
				userEmail = email
				c.Set("user_email", userEmail)
			}
			c.Set("role", claims["role"])
			logger.SystemLogger.Info().Msgf("DEBUG: AuthMiddleware set role=%v", claims["role"])

			// Also set in request context for logger.WithCtx compatibility
			ctx := c.Request.Context()
			if userID != 0 {
				ctx = context.WithValue(ctx, logger.CtxKeyUserID, userID)
			}
			if userEmail != "" {
				ctx = context.WithValue(ctx, logger.CtxKeyUserEmail, userEmail)
			}
			c.Request = c.Request.WithContext(ctx)

			if permissionsMaskFloat, ok := claims["permissions_mask"].(float64); ok {
				permissionsMask := uint64(permissionsMaskFloat)
				c.Set("permissions_mask", permissionsMask)
				logger.SystemLogger.Info().Msgf("DEBUG: AuthMiddleware set permissions_mask=%d", permissionsMask)
			} else {
				logger.SystemLogger.Info().Msgf("DEBUG: AuthMiddleware permissions_mask claim missing or invalid type: %T", claims["permissions_mask"])
				// Set empty mask to avoid "not found" error in guard
				c.Set("permissions_mask", uint64(0))
			}
		} else {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		c.Next()
	}
}
