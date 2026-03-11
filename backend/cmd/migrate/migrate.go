package main

import (
	"github.com/google/uuid"
	"github.com/hadi-projects/go-react-starter/config"
	customEntity "github.com/hadi-projects/go-react-starter/internal/entity"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"github.com/hadi-projects/go-react-starter/pkg/database"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(&cfg)

	db, err := database.NewMySQLConnection(&cfg)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	logger.SystemLogger.Info().Msg("Starting auto-migration...")

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.PasswordResetToken{},
		&entity.HttpLog{},
		&customEntity.Testsaja{},
		&customEntity.Produk{},
		&customEntity.Testdua{},
		// [GENERATOR_INSERT_MIGRATION]
	)

	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to auto-migrate database")
	}

	// Log this action to the new http_logs table
	logAction := &entity.HttpLog{
		RequestID:       uuid.New().String(),
		Method:          "SYSTEM",
		Path:            "database:migrate",
		ClientIP:        "127.0.0.1",
		UserAgent:       "Go-React-Starter/CLI",
		RequestHeaders:  "{}",
		RequestBody:     "{}",
		StatusCode:      200,
		ResponseHeaders: "{}",
		ResponseBody:    `{"message": "Auto-migration completed successfully"}`,
		Latency:         0,
		UserEmail:       "system@local",
	}
	db.Create(logAction)

	logger.SystemLogger.Info().Msg("Auto-migration completed successfully!")
}
