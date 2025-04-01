package routes

import (
	customMiddleware "goapi-starter/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, specify exact domains
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// Global middleware
	r.Use(corsMiddleware.Handler)
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
		r.Mount("/api/dummy-products", DummyProductRoutes())
	})

	return r
}
