package main

import (
	"cursor-experiment-1/internal/database"
	"cursor-experiment-1/internal/models"
	"cursor-experiment-1/internal/routes"
	"net/http"
)

func main() {
	// Initialize database
	database.InitDB()

	// Auto migrate the schema
	database.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.RefreshToken{})

	router := routes.SetupRouter()

	println("Server starting on :3000")
	http.ListenAndServe(":3000", router)
}
