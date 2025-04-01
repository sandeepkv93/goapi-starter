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

	// HandlerCounter counts handler executions
	HandlerCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_handler_executions_total",
			Help: "Total number of handler executions",
		},
		[]string{"handler", "status"},
	)

	// HandlerDuration measures handler execution duration
	HandlerDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "goapi_handler_duration_seconds",
			Help:    "Duration of handler executions in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "status"},
	)

	// HandlerErrors counts handler errors
	HandlerErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_handler_errors_total",
			Help: "Total number of handler errors",
		},
		[]string{"handler", "error_type"},
	)

	// BusinessOperations counts business logic operations
	BusinessOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_business_operations_total",
			Help: "Total number of business operations",
		},
		[]string{"operation", "result"},
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

// RecordHandlerExecution records metrics for a handler execution
func RecordHandlerExecution(handler string, status int, duration time.Duration) {
	HandlerCounter.WithLabelValues(handler, string(rune(status))).Inc()
	HandlerDuration.WithLabelValues(handler, string(rune(status))).Observe(duration.Seconds())
}

// RecordHandlerError records a handler error
func RecordHandlerError(handler, errorType string) {
	HandlerErrors.WithLabelValues(handler, errorType).Inc()
}

// RecordBusinessOperation records a business operation
func RecordBusinessOperation(operation, result string) {
	BusinessOperations.WithLabelValues(operation, result).Inc()
}
