package metrics

import (
	"strconv"
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

	// ErrorDetails provides more detailed error information
	ErrorDetails = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_error_details_total",
			Help: "Detailed breakdown of errors by type and reason",
		},
		[]string{"handler", "error_type", "error_reason"},
	)

	// CacheOperations counts cache operations
	CacheOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "type"},
	)

	// CacheResults counts cache hits/misses
	CacheResults = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapi_cache_results_total",
			Help: "Total number of cache hits/misses",
		},
		[]string{"result"},
	)

	// CacheDuration measures the time taken for cache operations
	CacheDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "goapi_cache_duration_seconds",
			Help:    "Duration of cache operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// CacheSize tracks the size of cached objects
	CacheSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "goapi_cache_size_bytes",
			Help:    "Size of cached objects in bytes",
			Buckets: []float64{64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536},
		},
		[]string{"key_prefix"},
	)
)

// RecordRequest records metrics for an HTTP request
func RecordRequest(method, path string, status int, duration time.Duration) {
	statusStr := strconv.Itoa(status)
	RequestCounter.WithLabelValues(method, path, statusStr).Inc()
	RequestDuration.WithLabelValues(method, path, statusStr).Observe(duration.Seconds())
}

// RecordDatabaseOperation records a database operation
func RecordDatabaseOperation(operation, entity string) {
	DatabaseOperations.WithLabelValues(operation, entity).Inc()
}

// RecordHandlerExecution records metrics for a handler execution
func RecordHandlerExecution(handler string, status int, duration time.Duration) {
	statusStr := strconv.Itoa(status)
	HandlerCounter.WithLabelValues(handler, statusStr).Inc()
	HandlerDuration.WithLabelValues(handler, statusStr).Observe(duration.Seconds())
}

// RecordHandlerError records a handler error
func RecordHandlerError(handler, errorType string) {
	HandlerErrors.WithLabelValues(handler, errorType).Inc()
}

// RecordDetailedError records a detailed error with reason
func RecordDetailedError(handler, errorType, reason string) {
	HandlerErrors.WithLabelValues(handler, errorType).Inc()
	ErrorDetails.WithLabelValues(handler, errorType, reason).Inc()
}

// RecordBusinessOperation records a business operation
func RecordBusinessOperation(operation, result string) {
	BusinessOperations.WithLabelValues(operation, result).Inc()
}

// RecordCacheOperation records a cache operation
func RecordCacheOperation(operation, opType string) {
	CacheOperations.WithLabelValues(operation, opType).Inc()
}

// RecordCacheResult records a cache hit or miss
func RecordCacheResult(result string) {
	CacheResults.WithLabelValues(result).Inc()
}

// RecordCacheDuration records the duration of a cache operation
func RecordCacheDuration(operation string, duration time.Duration) {
	CacheDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordCacheSize records the size of a cached object
func RecordCacheSize(keyPrefix string, size int) {
	CacheSize.WithLabelValues(keyPrefix).Observe(float64(size))
}
