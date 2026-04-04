package router

import (
	"github.com/gin-gonic/gin"
	customHandler "github.com/hadi-projects/go-react-starter/internal/handler"
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
	httpLogHandler handler.HttpLogHandler,
	systemLogHandler handler.SystemLogHandler,
	auditLogHandler handler.AuditLogHandler,
	generatorHandler handler.GeneratorHandler,
	produkHandler customHandler.ProdukHandler,
	storageHandler customHandler.StorageHandler,
	healthHandler handler.HealthHandler,
	settingHandler handler.SettingHandler,
	apiKeyHandler handler.ApiKeyHandler,
	permGuard *middleware.PermissionGuard,
	// [GENERATOR_INSERT_HANDLER_PARAM]
) {
	// Health and Status
	health := v1.Group("/health")
	{
		health.GET("/status", healthHandler.GetStatus)
	}

	// Module Generator
	generator := v1.Group("/generator")
	generator.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		generator.POST("", permGuard.Check("create-module"), generatorHandler.Generate)
	}
	// Storage routes (authenticated)
	storageGroup := v1.Group("/storage")
	storageGroup.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		storageGroup.POST("/upload", permGuard.Check("upload-file"), storageHandler.Upload)
		storageGroup.GET("", permGuard.Check("get-file"), storageHandler.GetFiles)
		storageGroup.GET("/:id", permGuard.Check("get-file"), storageHandler.GetFileByID)
		storageGroup.DELETE("/:id", permGuard.Check("delete-file"), storageHandler.DeleteFile)
		storageGroup.GET("/:id/download", permGuard.Check("get-file"), storageHandler.DownloadFile)
		storageGroup.POST("/:id/share", permGuard.Check("share-file"), storageHandler.CreateShareLink)
		storageGroup.GET("/:id/shares", permGuard.Check("share-file"), storageHandler.GetShareLinks)
		storageGroup.PUT("/shares/:shareId", permGuard.Check("share-file"), storageHandler.UpdateShareLink)
		storageGroup.DELETE("/shares/:shareId", permGuard.Check("share-file"), storageHandler.RevokeShareLink)
		storageGroup.GET("/shares/:shareId/logs", permGuard.Check("share-file"), storageHandler.GetShareLinkLogs)
	}

	// Public share routes (no auth required)
	publicGroup := v1.Group("/public")
	{
		publicGroup.GET("/share/:token", storageHandler.PublicFileInfo)
		publicGroup.GET("/share/:token/download", storageHandler.PublicDownload)
		publicGroup.GET("/storage/:id", storageHandler.PublicSystemFile)
		publicGroup.GET("/settings/:category", settingHandler.GetPublicByCategory)
	}

	produk := v1.Group("/produk")
	produk.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		produk.POST("", produkHandler.Create)
		produk.GET("", produkHandler.GetAll)
		produk.GET("/:id", produkHandler.GetByID)
		produk.PUT("/:id", produkHandler.Update)
		produk.DELETE("/:id", produkHandler.Delete)
		produk.GET("/export", produkHandler.Export)
	}
	// [GENERATOR_INSERT_GROUP]
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", middleware.AuthMiddleware(r.config.JWT.Secret), authHandler.Logout)
		auth.POST("/register", userHandler.Register)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh", authHandler.RefreshToken)
		// 2FA routes
		auth.POST("/2fa/verify", authHandler.Verify2FA) // Public: no JWT needed
		auth.POST("/2fa/enroll", middleware.AuthMiddleware(r.config.JWT.Secret), authHandler.Enroll2FA)
		auth.POST("/2fa/confirm", middleware.AuthMiddleware(r.config.JWT.Secret), authHandler.Confirm2FA)
		auth.DELETE("/2fa/disable", middleware.AuthMiddleware(r.config.JWT.Secret), authHandler.Disable2FA)
		auth.POST("/2fa/reset-request", authHandler.Request2FAReset) // Public: needs temp token
		auth.POST("/2fa/reset-confirm", authHandler.Confirm2FAReset) // Public: needs email format token
	}

	logs := v1.Group("/logs")
	logs.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		// Internal permission check is handled inside GetLogs
		logs.GET("", logHandler.GetLogs)
		logs.GET("/export", logHandler.Export) // Combined logs export
		logs.GET("/http", permGuard.Check("get-http-log"), httpLogHandler.GetAll)
		logs.GET("/http/export", permGuard.Check("get-http-log"), httpLogHandler.Export)
		logs.GET("/system", permGuard.Check("get-http-log"), systemLogHandler.GetAll) // Use same permission for now
		logs.GET("/system/export", permGuard.Check("get-http-log"), systemLogHandler.Export)
		logs.GET("/audit", permGuard.Check("get-http-log"), auditLogHandler.GetAll) // Use same permission for now
		logs.GET("/audit/export", permGuard.Check("get-http-log"), auditLogHandler.Export)
	}

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		// User can access their own profile
		users.GET("/me", userHandler.Me)

		// Admin only for CRUD
		users.POST("", permGuard.Check("create-user"), userHandler.Create)
		users.GET("", permGuard.Check("get-user"), userHandler.GetAll)
		users.GET("/export", permGuard.Check("get-user"), userHandler.Export)
		users.PUT("/:id", permGuard.Check("edit-user"), userHandler.Update)
		users.DELETE("/:id", permGuard.Check("delete-user"), userHandler.Delete)
	}

	permissions := v1.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	permissions.Use(permGuard.Check("get-permission")) // Assuming admin role has this
	{
		permissions.POST("", permGuard.Check("create-permission"), permissionHandler.Create)
		permissions.GET("", permGuard.Check("get-permission"), permissionHandler.GetAll)
		permissions.GET("/export", permGuard.Check("get-permission"), permissionHandler.Export)
		permissions.PUT("/:id", permGuard.Check("edit-permission"), permissionHandler.Update)
		permissions.DELETE("/:id", permGuard.Check("delete-permission"), permissionHandler.Delete)
	}

	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	roles.Use(permGuard.Check("get-role"))
	{
		roles.POST("", permGuard.Check("create-role"), roleHandler.Create)
		roles.GET("", permGuard.Check("get-role"), roleHandler.GetAll)
		roles.GET("/export", permGuard.Check("get-role"), roleHandler.Export)
		roles.GET("/:id", permGuard.Check("get-role"), roleHandler.GetByID)
		roles.PUT("/:id", permGuard.Check("edit-role"), roleHandler.Update)
		roles.DELETE("/:id", permGuard.Check("delete-role"), roleHandler.Delete)
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
		cache.DELETE("/clear", permGuard.Check("manage-cache"), cacheHandler.ClearAll)
		cache.GET("/status", permGuard.Check("manage-cache"), cacheHandler.GetStatus)
	}

	// Settings management
	settings := v1.Group("/settings")
	settings.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		settings.GET("/:category", permGuard.Check("get-setting"), settingHandler.GetByCategory)
		settings.PUT("", permGuard.Check("edit-setting"), settingHandler.BulkUpdate)
	}

	// API Keys management
	apiKeys := v1.Group("/apikeys")
	apiKeys.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		apiKeys.POST("", permGuard.Check("create-api-key"), apiKeyHandler.Create)
		apiKeys.GET("", permGuard.Check("get-api-key"), apiKeyHandler.GetAll)
		apiKeys.DELETE("/:id", permGuard.Check("delete-api-key"), apiKeyHandler.Delete)
	}
}
