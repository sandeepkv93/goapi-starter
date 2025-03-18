package routes

import (
	customMiddleware "goapi-starter/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(customMiddleware.LoggingMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Mount("/api/auth", AuthRoutes())

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)
		r.Mount("/api/users", UserRoutes())
		r.Mount("/api/products", ProductRoutes())
	})

	return r
}
