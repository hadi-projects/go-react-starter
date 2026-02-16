package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/internal/handler"
	"github.com/hadi-projects/go-react-starter/internal/middleware"
	"github.com/hadi-projects/go-react-starter/internal/repository"
	"github.com/hadi-projects/go-react-starter/internal/service"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/database"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type Router struct {
	config *config.Config
	cache  cache.CacheService
}

func NewRouter(config *config.Config, cache cache.CacheService) *Router {
	return &Router{
		config: config,
		cache:  cache,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	if r.config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	db, err := database.NewMySQLConnection(r.config)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(r.config))
	router.Use(middleware.RequestLogger())
	router.Use(middleware.RequestCancellation(time.Duration(r.config.Security.RequestTimeOut) * time.Second))
	// router.Use(middleware.APIKeyMiddleware(r.config.Security.APIKey)) // Removed global application
	router.Use(middleware.RateLimiter(r.config.RateLimit.Rps, r.config.RateLimit.Burst))
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.XSSProtection())

	// Repositories
	userRepo := repository.NewUserRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, r.config)
	userService := service.NewUserService(userRepo, r.config, r.cache)
	permissionService := service.NewPermissionService(permissionRepo, r.cache)
	roleService := service.NewRoleService(roleRepo, r.cache)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	roleHandler := handler.NewRoleHandler(roleService)
	cacheHandler := handler.NewCacheHandler(r.cache)

	v1 := router.Group("/api/v1")
	{
		r.setupPrivateRoutes(v1, authHandler, userHandler, permissionHandler, roleHandler, cacheHandler)
	}

	logger.SystemLogger.Info().Str("port", r.config.App.Port).Msg("Server running")
	return router
}

func (r *Router) Run() {
	srv := &http.Server{
		Addr:           ":" + r.config.App.Port,
		Handler:        r.SetupRouter(),
		ReadTimeout:    time.Duration(r.config.Security.RequestTimeOut) * time.Second,
		WriteTimeout:   time.Duration(r.config.Security.RequestTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		logger.SystemLogger.Info().Str("port", r.config.App.Port).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.SystemLogger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	fmt.Println("Server exited successfully")

}
