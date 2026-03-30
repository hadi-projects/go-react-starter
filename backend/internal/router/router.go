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
	customHandler "github.com/hadi-projects/go-react-starter/internal/handler"
	handler "github.com/hadi-projects/go-react-starter/internal/handler/default"
	"github.com/hadi-projects/go-react-starter/internal/middleware"
	customRepository "github.com/hadi-projects/go-react-starter/internal/repository"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	customService "github.com/hadi-projects/go-react-starter/internal/service"
	service "github.com/hadi-projects/go-react-starter/internal/service/default"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/kafka"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/mailer"
	"gorm.io/gorm"
)

type Router struct {
	config        *config.Config
	db            *gorm.DB
	cache         cache.CacheService
	kafkaProducer kafka.Producer
	mailer        mailer.Mailer
}

func NewRouter(config *config.Config, db *gorm.DB, cache cache.CacheService, kafkaProducer kafka.Producer, mailer mailer.Mailer) *Router {
	return &Router{
		config:        config,
		db:            db,
		cache:         cache,
		kafkaProducer: kafkaProducer,
		mailer:        mailer,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	if r.config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Use db from Router struct
	db := r.db

	// Repositories initializations for middleware
	httpLogRepo := repository.NewHttpLogRepository(db)
	systemLogRepo := repository.NewSystemLogRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(r.config))
	router.Use(middleware.RequestLogger(r.config, httpLogRepo))
	router.Use(middleware.RequestCancellation(time.Duration(r.config.Security.RequestTimeOut) * time.Second))
	router.Use(middleware.RateLimiter(r.config.RateLimit.Rps, r.config.RateLimit.Burst))
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.XSSProtection())

	// Repositories
	userRepo := repository.NewUserRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	produkRepo := customRepository.NewProdukRepository(db)
	// [GENERATOR_INSERT_REPOSITORY]

	// Services
	authService := service.NewAuthService(userRepo, tokenRepo, r.kafkaProducer, r.mailer, r.config, r.cache)
	userService := service.NewUserService(userRepo, r.config, r.cache)
	permissionService := service.NewPermissionService(permissionRepo, r.cache)
	roleService := service.NewRoleService(roleRepo, r.cache)
	logService := service.NewLogService(r.config)
	statisticsService := service.NewStatisticsService(db)
	httpLogService := service.NewHttpLogService(httpLogRepo, r.cache)
	systemLogService := service.NewSystemLogService(systemLogRepo, r.cache)
	auditLogService := service.NewAuditLogService(auditLogRepo, r.cache)
	produkService := customService.NewProdukService(produkRepo, r.cache)
	// [GENERATOR_INSERT_SERVICE]

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	roleHandler := handler.NewRoleHandler(roleService)
	logHandler := handler.NewLogHandler(logService)
	cacheHandler := handler.NewCacheHandler(r.cache)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)
	httpLogHandler := handler.NewHttpLogHandler(httpLogService)
	systemLogHandler := handler.NewSystemLogHandler(systemLogService)
	auditLogHandler := handler.NewAuditLogHandler(auditLogService)
	generatorHandler := handler.NewGeneratorHandler(".", db)
	produkHandler := customHandler.NewProdukHandler(produkService)
	healthHandler := handler.NewHealthHandler(r.cache, r.kafkaProducer)
	// [GENERATOR_INSERT_HANDLER]

	v1 := router.Group("/api/v1")
	{
		r.setupPrivateRoutes(v1, authHandler, userHandler, permissionHandler, roleHandler, logHandler, cacheHandler, statisticsHandler, httpLogHandler, systemLogHandler, auditLogHandler,

			generatorHandler,
			produkHandler,
			healthHandler,
			middleware.NewPermissionGuard(r.cache, permissionRepo),
		// [GENERATOR_INSERT_HANDLER_PARAM]
		)
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
