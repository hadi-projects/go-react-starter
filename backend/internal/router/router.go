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
	"github.com/hadi-projects/go-react-starter/internal/middleware"
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

	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(r.config))
	router.Use(middleware.RequestLogger())
	router.Use(middleware.RequestCancellation(time.Duration(r.config.Security.RequestTimeOut) * time.Second))
	router.Use(middleware.APIKeyMiddleware(r.config.Security.APIKey))
	router.Use(middleware.RateLimiter(r.config.RateLimit.Rps, r.config.RateLimit.Burst))
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.XSSProtection())

	r.setupPublicRuotes(router)

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
