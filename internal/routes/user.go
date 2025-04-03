package routes

import (
	"goapi-starter/internal/handlers"
	"goapi-starter/internal/middleware"
	"goapi-starter/internal/utils"

	"github.com/go-chi/chi/v5"
)

func UserRoutes() chi.Router {
	r := chi.NewRouter()

	// Apply auth middleware to all routes in this group
	r.Use(middleware.AuthMiddleware)

	// Protected routes
	r.Get("/profile", utils.InstrumentHandler("GetProfile", handlers.GetProfile))

	return r
}
