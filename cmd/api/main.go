package main

import (
	"goapi-starter/internal/config"
	"goapi-starter/internal/database"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/models"
	"goapi-starter/internal/routes"
	"net/http"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDB()

	// Auto migrate the schema
	database.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.RefreshToken{})

	router := routes.SetupRouter()

	port := config.AppConfig.Server.Port
	logger.Info().Msgf("Server starting on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
