package main

import (
	"log"

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
		log.Fatal("Failed to connect to database: ", err)
	}

	cacheService, err := cache.NewRedisCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}
	defer cacheService.Close()

	router := router.NewRouter(&cfg, cacheService)
	router.SetupRouter()
	router.Run()
}
