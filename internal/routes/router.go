package routes

import (
	"goapi-starter/internal/handlers"
	customMiddleware "goapi-starter/internal/middleware"
	"goapi-starter/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Add correlation middleware early in the chain
	r.Use(customMiddleware.CorrelationMiddleware)
	r.Use(customMiddleware.LoggingMiddleware)

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
	r.Use(customMiddleware.PrometheusMiddleware)
	r.Use(middleware.Recoverer)

	// Apply IP-based rate limiting to all routes
	r.Use(customMiddleware.IPRateLimitMiddleware)

	// Monitoring endpoints (health & metrics)
	// These only need basic rate limiting
	r.Group(func(r chi.Router) {
		r.Get("/health", utils.InstrumentHandler("HealthCheck", handlers.HealthCheck))
		r.Handle("/metrics", promhttp.Handler())
	})

	// Auth routes with stricter rate limiting
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthRateLimitMiddleware)
		r.Mount("/api/auth", AuthRoutes())
	})

	// Protected routes with user-based rate limiting
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)
		r.Use(customMiddleware.UserRateLimitMiddleware)
		r.Mount("/api/dummy-products", DummyProductRoutes())

		// User routes
		r.Mount("/api/user", UserRoutes())

		// Logout route
		r.Post("/api/auth/logout", utils.InstrumentHandler("Logout", handlers.Logout))
	})

	return r
}
