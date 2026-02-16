package main

import (
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/internal/router"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/database"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(&cfg)

	_, err := database.NewMySQLConnection(&cfg)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	cacheService, err := cache.NewRedisCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer cacheService.Close()

	router := router.NewRouter(&cfg, cacheService)
	router.SetupRouter()
	router.Run()
}
