package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestCounter counts the total number of HTTP requests
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// RequestDuration measures the duration of HTTP requests
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "goapi_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// ActiveRequests tracks the number of currently active requests
	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goapi_http_requests_active",
			Help: "Number of currently active HTTP requests",
		},
	)

	// DatabaseOperations counts database operations
	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "entity"},
	)
)

// RecordRequest records metrics for an HTTP request
func RecordRequest(method, path string, status int, duration time.Duration) {
	RequestCounter.WithLabelValues(method, path, string(rune(status))).Inc()
	RequestDuration.WithLabelValues(method, path, string(rune(status))).Observe(duration.Seconds())
}

// RecordDatabaseOperation records a database operation
func RecordDatabaseOperation(operation, entity string) {
	DatabaseOperations.WithLabelValues(operation, entity).Inc()
}
