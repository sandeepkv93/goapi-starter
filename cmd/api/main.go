package main

import (
	"cursor-experiment-1/internal/config"
	"cursor-experiment-1/internal/database"
	"cursor-experiment-1/internal/models"
	"cursor-experiment-1/internal/routes"
	"fmt"
	"net/http"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDB()

	// Auto migrate the schema
	database.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.RefreshToken{})

	router := routes.SetupRouter()

	port := config.AppConfig.Server.Port
	fmt.Printf("Server starting on :%s\n", port)
	http.ListenAndServe(":"+port, router)
}
