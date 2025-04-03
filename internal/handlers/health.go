package handlers

import (
	"goapi-starter/internal/database"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/utils"
	"net/http"
)

// HealthCheck handles health check requests
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	logger.Info().Msg("Checking database connection")
	sqlDB, err := database.DB.DB()
	if err != nil {
		logger.Error().Err(err).Msg("Database connection error")
		utils.RespondWithError(w, r, http.StatusServiceUnavailable, "Database connection error")
		return
	}
	logger.Info().Msg("Database connection successful")

	// Ping the database
	logger.Info().Msg("Pinging database")
	if err := sqlDB.Ping(); err != nil {
		logger.Error().Err(err).Msg("Database ping failed")
		utils.RespondWithError(w, r, http.StatusServiceUnavailable, "Database ping failed")
		return
	}
	logger.Info().Msg("Database ping successful")

	// All checks passed
	logger.Info().Msg("Service is healthy")
	utils.RespondWithJSON(w, r, http.StatusOK, utils.SuccessResponse{
		Message: "Service is healthy",
		Data: map[string]string{
			"status": "UP",
		},
	})
}
