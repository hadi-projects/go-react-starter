package main

import (
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/internal/router"
	"github.com/hadi-projects/go-react-starter/pkg/cache"
	"github.com/hadi-projects/go-react-starter/pkg/database"
	"github.com/hadi-projects/go-react-starter/pkg/kafka"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/mailer"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(&cfg)

	db, err := database.NewMySQLConnection(&cfg)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	cacheService, err := cache.NewRedisCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer cacheService.Close()

	kafkaProducer, err := kafka.NewProducer(&cfg)
	if err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Failed to create Kafka producer")
		kafkaProducer = nil
	} else {
		defer kafkaProducer.Close()
	}

	mailService := mailer.NewMailer(&cfg)

	router := router.NewRouter(&cfg, db, cacheService, kafkaProducer, mailService)
	router.SetupRouter()
	router.Run()
}
