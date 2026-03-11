package main

import (
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/database"
	"github.com/hadi-projects/go-react-starter/pkg/database/seeder"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/google/uuid"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(&cfg)

	db, err := database.NewMySQLConnection(&cfg)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	seeder.SeedRole(db)
	seeder.SeedUser(db, cfg.Security.BCryptCost)

	// Log this action to the new http_logs table
	logAction := &entity.HttpLog{
		RequestID:       uuid.New().String(),
		Method:          "SYSTEM",
		Path:            "database:seed",
		ClientIP:        "127.0.0.1",
		UserAgent:       "Go-React-Starter/CLI",
		RequestHeaders:  "{}",
		RequestBody:     "{}",
		StatusCode:      200,
		ResponseHeaders: "{}",
		ResponseBody:    `{"message": "Database seeder completed successfully"}`,
		Latency:         0,
		UserEmail:       "system@local",
	}
	db.Create(logAction)
}
