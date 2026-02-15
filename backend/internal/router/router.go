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
	if r.config.APPEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.RateLimiter(r.config.RateLimitRps, r.config.RateLimitBurst))
	router.Use(middleware.CORS(r.config))
	router.Use(middleware.RequestLogger())
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.XSSProtection())

	r.setupPublicRuotes(router)

	return router
}

func (r *Router) Run() {
	app := &http.Server{
		Addr:           fmt.Sprintf(":%s", r.config.AppPort),
		Handler:        r.SetupRouter(),
		ReadTimeout:    time.Duration(r.config.RequestTimeOut),
		WriteTimeout:   time.Duration(r.config.RequestTimeOut),
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		fmt.Printf("Server running on port :%s", r.config.AppPort)
		if err := app.ListenAndServe(); err != nil {
			log.Fatal("Server failed start: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown")
	}

	fmt.Println("Server exited successfully")

}
