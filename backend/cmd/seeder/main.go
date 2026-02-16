package main

import (
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/database"
	"github.com/hadi-projects/go-react-starter/pkg/database/seeder"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(&cfg)

	db, err := database.NewMySQLConnection(&cfg)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	seeder.SeedUser(db, cfg.Security.BCryptCost)
	seeder.SeedRole(db)
}
