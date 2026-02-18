package main

import (
	"github.com/hadi-projects/go-react-starter/config"
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
		// [GENERATOR_INSERT_MIGRATION]
	)

	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to auto-migrate database")
	}

	logger.SystemLogger.Info().Msg("Auto-migration completed successfully!")
}
