package middleware

import (
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"net/http"
	"time"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Track active requests
		metrics.ActiveRequests.Inc()
		defer metrics.ActiveRequests.Dec()

		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("Recording metrics for request")

		// Use our existing responseWriter wrapper
		wrapped := wrapResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Record metrics after the request is complete
		metrics.RecordRequest(
			r.Method,
			r.URL.Path,
			wrapped.status,
			duration,
		)

		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapped.status).
			Dur("duration", duration).
			Msg("Recorded metrics for request")
	})
}
