package routes

import (
	"goapi-starter/internal/handlers"
	"goapi-starter/internal/utils"

	"github.com/go-chi/chi/v5"
)

// HealthRoutes sets up health check routes
func HealthRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", utils.InstrumentHandler("HealthCheck", handlers.HealthCheck))

	return router
}
