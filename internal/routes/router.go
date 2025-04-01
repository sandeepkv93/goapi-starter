package routes

import (
	customMiddleware "goapi-starter/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(customMiddleware.LoggingMiddleware)
	r.Use(customMiddleware.PrometheusMiddleware)
	r.Use(middleware.Recoverer)

	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

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
