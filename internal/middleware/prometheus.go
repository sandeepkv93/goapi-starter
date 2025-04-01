package middleware

import (
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

		// Use our existing responseWriter wrapper
		wrapped := wrapResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Record metrics after the request is complete
		metrics.RecordRequest(
			r.Method,
			r.URL.Path,
			wrapped.status,
			time.Since(start),
		)
	})
}
