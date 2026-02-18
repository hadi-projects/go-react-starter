package router

import (
	"github.com/gin-gonic/gin"
	handler "github.com/hadi-projects/go-react-starter/internal/handler/default"
	"github.com/hadi-projects/go-react-starter/internal/middleware"
)

func (r *Router) setupPrivateRoutes(
	v1 *gin.RouterGroup,
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
	permissionHandler handler.PermissionHandler,
	roleHandler handler.RoleHandler,
	logHandler handler.LogHandler,
	cacheHandler handler.CacheHandler,
	statisticsHandler handler.StatisticsHandler,
	generatorHandler handler.GeneratorHandler,
	// [GENERATOR_INSERT_HANDLER_PARAM]
) {
	// Module Generator
	generator := v1.Group("/generator")
	generator.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		generator.POST("", middleware.PermissionGuard("create-module"), generatorHandler.Generate)
	}
	// [GENERATOR_INSERT_GROUP]
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", userHandler.Register)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	logs := v1.Group("/logs")
	logs.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		// Internal permission check is handled inside GetLogs
		logs.GET("", logHandler.GetLogs)
	}

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		// User can access their own profile
		users.GET("/me", middleware.PermissionGuard("get-profile"), userHandler.Me)

		// Admin only for CRUD
		users.POST("", middleware.PermissionGuard("create-user"), userHandler.Create)
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

	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	roles.Use(middleware.PermissionGuard("get-role"))
	{
		roles.POST("", middleware.PermissionGuard("create-role"), roleHandler.Create)
		roles.GET("", middleware.PermissionGuard("get-role"), roleHandler.GetAll)
		roles.GET("/:id", middleware.PermissionGuard("get-role"), roleHandler.GetByID)
		roles.PUT("/:id", middleware.PermissionGuard("edit-role"), roleHandler.Update)
		roles.DELETE("/:id", middleware.PermissionGuard("delete-role"), roleHandler.Delete)
	}

	// Statistics
	statistics := v1.Group("/statistics")
	statistics.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		statistics.GET("/dashboard", statisticsHandler.GetDashboardStats)
	}

	// Cache management
	cache := v1.Group("/cache")
	cache.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		cache.DELETE("/clear", middleware.PermissionGuard("manage-cache"), cacheHandler.ClearAll)
		cache.GET("/status", middleware.PermissionGuard("manage-cache"), cacheHandler.GetStatus)
	}
}
