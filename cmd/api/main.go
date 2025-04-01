package main

import (
	"context"
	"goapi-starter/internal/cache"
	"goapi-starter/internal/config"
	"goapi-starter/internal/database"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/models"
	"goapi-starter/internal/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Info().Msg("Starting application")

	// Load configuration
	logger.Info().Msg("Loading configuration")
	config.LoadConfig()

	// Initialize database
	logger.Info().Msg("Initializing database connection")
	database.InitDB()

	// Initialize Redis
	logger.Info().Msg("Initializing Redis connection")
	cache.InitRedis()

	// Auto migrate the schema
	logger.Info().Msg("Running database migrations")
	if err := database.DB.AutoMigrate(&models.User{}, &models.DummyProduct{}, &models.RefreshToken{}); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	logger.Info().Msg("Database migrations completed successfully")

	// Setup router
	logger.Info().Msg("Setting up HTTP routes")
	router := routes.SetupRouter()

	// Start server
	port := config.AppConfig.Server.Port
	logger.Info().
		Str("port", port).
		Msg("Server starting")
	logger.Info().Msg("Metrics available at /metrics")

	// Setup graceful shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited gracefully")
}
