package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/internal/handler"
	"github.com/hadi-projects/go-react-starter/internal/middleware"
)

func (r *Router) setupPrivateRoutes(
	v1 *gin.RouterGroup,
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
	permissionHandler handler.PermissionHandler,
) {
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", userHandler.Register)
	}

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		// User can access their own profile
		users.GET("/me", middleware.PermissionGuard("get-profile"), userHandler.Me)

		// Admin only for CRUD
		users.GET("", middleware.PermissionGuard("get-user"), userHandler.GetAll)
		users.PUT("/:id", middleware.PermissionGuard("edit-user"), userHandler.Update)
		users.DELETE("/:id", middleware.PermissionGuard("delete-user"), userHandler.Delete)
	}

	permissions := v1.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	permissions.Use(middleware.PermissionGuard("get-permission")) // Assuming admin role has this
	{
		permissions.POST("", middleware.PermissionGuard("create-permission"), permissionHandler.Create)
		permissions.GET("", middleware.PermissionGuard("get-permission"), permissionHandler.GetAll)
		permissions.PUT("/:id", middleware.PermissionGuard("edit-permission"), permissionHandler.Update)
		permissions.DELETE("/:id", middleware.PermissionGuard("delete-permission"), permissionHandler.Delete)
	}
}
