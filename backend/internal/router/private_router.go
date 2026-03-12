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
	testsajaHandler customHandler.TestsajaHandler,
	produkHandler customHandler.ProdukHandler,
	healthHandler handler.HealthHandler,
	testduaHandler customHandler.TestduaHandler,
	cookHandler customHandler.CookHandler,
	adminHandler customHandler.AdminHandler,
	mainnnHandler customHandler.MainnnHandler,
	wisudaHandler customHandler.WisudaHandler,
		arsipHandler customHandler.ArsipHandler,
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
		generator.POST("", middleware.PermissionGuard("create-module"), generatorHandler.Generate)
	}
	testsaja := v1.Group("/testsaja")
	testsaja.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		testsaja.POST("", testsajaHandler.Create)
		testsaja.GET("", testsajaHandler.GetAll)
		testsaja.GET("/:id", testsajaHandler.GetByID)
		testsaja.PUT("/:id", testsajaHandler.Update)
		testsaja.DELETE("/:id", testsajaHandler.Delete)
		testsaja.GET("/export", testsajaHandler.Export)
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
	testdua := v1.Group("/testdua")
	testdua.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		testdua.POST("", testduaHandler.Create)
		testdua.GET("", testduaHandler.GetAll)
		testdua.GET("/:id", testduaHandler.GetByID)
		testdua.PUT("/:id", testduaHandler.Update)
		testdua.DELETE("/:id", testduaHandler.Delete)
		testdua.GET("/export", testduaHandler.Export)
	}
	cook := v1.Group("/cook")
	cook.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		cook.POST("", cookHandler.Create)
		cook.GET("", cookHandler.GetAll)
		cook.GET("/:id", cookHandler.GetByID)
		cook.PUT("/:id", cookHandler.Update)
		cook.DELETE("/:id", cookHandler.Delete)
		cook.GET("/export", cookHandler.Export)
	}
	admin := v1.Group("/admin/module") // Use /admin/module to avoid conflict with /admin/*
	admin.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		admin.POST("", adminHandler.Create)
		admin.GET("", adminHandler.GetAll)
		admin.GET("/:id", adminHandler.GetByID)
		admin.PUT("/:id", adminHandler.Update)
		admin.DELETE("/:id", adminHandler.Delete)
		admin.GET("/export", adminHandler.Export)
	}
	mainnn := v1.Group("/mainnn")
	mainnn.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		mainnn.POST("", mainnnHandler.Create)
		mainnn.GET("", mainnnHandler.GetAll)
		mainnn.GET("/:id", mainnnHandler.GetByID)
		mainnn.PUT("/:id", mainnnHandler.Update)
		mainnn.DELETE("/:id", mainnnHandler.Delete)
	}
		wisuda := v1.Group("/wisuda")
	wisuda.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		wisuda.POST("", wisudaHandler.Create)
		wisuda.GET("", wisudaHandler.GetAll)
		wisuda.GET("/export", wisudaHandler.Export)
		wisuda.GET("/:id", wisudaHandler.GetByID)
		wisuda.PUT("/:id", wisudaHandler.Update)
		wisuda.DELETE("/:id", wisudaHandler.Delete)
	}
		arsip := v1.Group("/arsip")
	arsip.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		arsip.POST("", arsipHandler.Create)
		arsip.GET("", arsipHandler.GetAll)
		arsip.GET("/export", arsipHandler.Export)
		arsip.GET("/:id", arsipHandler.GetByID)
		arsip.PUT("/:id", arsipHandler.Update)
		arsip.DELETE("/:id", arsipHandler.Delete)
	}
	// [GENERATOR_INSERT_GROUP]
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", middleware.AuthMiddleware(r.config.JWT.Secret), authHandler.Logout)
		auth.POST("/register", userHandler.Register)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	logs := v1.Group("/logs")
	logs.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		// Internal permission check is handled inside GetLogs
		logs.GET("", logHandler.GetLogs)
		logs.GET("/export", logHandler.Export) // Combined logs export
		logs.GET("/http", middleware.PermissionGuard("get-http-log"), httpLogHandler.GetAll)
		logs.GET("/http/export", middleware.PermissionGuard("get-http-log"), httpLogHandler.Export)
		logs.GET("/system", middleware.PermissionGuard("get-http-log"), systemLogHandler.GetAll) // Use same permission for now
		logs.GET("/system/export", middleware.PermissionGuard("get-http-log"), systemLogHandler.Export)
		logs.GET("/audit", middleware.PermissionGuard("get-http-log"), auditLogHandler.GetAll) // Use same permission for now
		logs.GET("/audit/export", middleware.PermissionGuard("get-http-log"), auditLogHandler.Export)
	}

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	{
		// User can access their own profile
		users.GET("/me", userHandler.Me)

		// Admin only for CRUD
		users.POST("", middleware.PermissionGuard("create-user"), userHandler.Create)
		users.GET("", middleware.PermissionGuard("get-user"), userHandler.GetAll)
		users.GET("/export", middleware.PermissionGuard("get-user"), userHandler.Export)
		users.PUT("/:id", middleware.PermissionGuard("edit-user"), userHandler.Update)
		users.DELETE("/:id", middleware.PermissionGuard("delete-user"), userHandler.Delete)
	}

	permissions := v1.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	permissions.Use(middleware.PermissionGuard("get-permission")) // Assuming admin role has this
	{
		permissions.POST("", middleware.PermissionGuard("create-permission"), permissionHandler.Create)
		permissions.GET("", middleware.PermissionGuard("get-permission"), permissionHandler.GetAll)
		permissions.GET("/export", middleware.PermissionGuard("get-permission"), permissionHandler.Export)
		permissions.PUT("/:id", middleware.PermissionGuard("edit-permission"), permissionHandler.Update)
		permissions.DELETE("/:id", middleware.PermissionGuard("delete-permission"), permissionHandler.Delete)
	}

	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
	roles.Use(middleware.PermissionGuard("get-role"))
	{
		roles.POST("", middleware.PermissionGuard("create-role"), roleHandler.Create)
		roles.GET("", middleware.PermissionGuard("get-role"), roleHandler.GetAll)
		roles.GET("/export", middleware.PermissionGuard("get-role"), roleHandler.Export)
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
