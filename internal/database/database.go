package database

import (
	"fmt"
	"goapi-starter/internal/config"
	"goapi-starter/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	logger.Debug().
		Str("host", config.AppConfig.Database.Host).
		Str("port", config.AppConfig.Database.Port).
		Str("dbname", config.AppConfig.Database.DBName).
		Str("user", config.AppConfig.Database.User).
		Str("sslmode", config.AppConfig.Database.SSLMode).
		Msg("Initializing database connection")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
	)

	var err error
	// Configure GORM without custom logger
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).
			Str("host", config.AppConfig.Database.Host).
			Str("port", config.AppConfig.Database.Port).
			Str("dbname", config.AppConfig.Database.DBName).
			Msg("Failed to connect to database")
	}

	logger.Info().
		Str("host", config.AppConfig.Database.Host).
		Str("port", config.AppConfig.Database.Port).
		Str("dbname", config.AppConfig.Database.DBName).
		Msg("Successfully connected to database")

	// Add metrics callbacks
	AddMetricsCallbacks()
}
