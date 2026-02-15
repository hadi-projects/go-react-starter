package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/internal/router"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type Application struct {
	Config *config.Config
	Server *http.Server
	Router *gin.Engine
}

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(&cfg)

	router := router.NewRouter(&cfg)
	router.SetupRouter()
	router.Run()
}
