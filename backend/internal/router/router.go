package router

import (
	"context"
	"fmt"
	"log"
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
	"github.com/hadi-projects/go-react-starter/pkg/database"
)

type Router struct {
	config *config.Config
}

func NewRouter(config *config.Config) *Router {
	return &Router{
		config: config,
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
		log.Fatal("Failed to connect to database: ", err)
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

	// Services
	authService := service.NewAuthService(userRepo, r.config)
	userService := service.NewUserService(userRepo, r.config)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", userHandler.Register)
		}

		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(r.config.JWT.Secret))
		{
			// User can access their own profile
			users.GET("/me", middleware.RoleGuard("user", "admin"), userHandler.Me)

			// Admin only for CRUD
			users.GET("", middleware.RoleGuard("admin"), userHandler.GetAll)
			users.PUT("/:id", middleware.RoleGuard("admin"), userHandler.Update)
			users.DELETE("/:id", middleware.RoleGuard("admin"), userHandler.Delete)
		}
	}

	log.Printf("Server running on port %s", r.config.App.Port)
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
		fmt.Printf("Server running on port :%s", r.config.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown")
	}

	fmt.Println("Server exited successfully")

}
