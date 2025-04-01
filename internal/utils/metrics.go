package utils

import (
	"goapi-starter/internal/metrics"
	"net/http"
	"time"
)

// InstrumentHandler wraps an HTTP handler with metrics instrumentation
func InstrumentHandler(handlerName string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := NewResponseWriter(w)

		// Execute the handler
		next(rw, r)

		// Record metrics
		duration := time.Since(start)
		metrics.RecordHandlerExecution(handlerName, rw.StatusCode(), duration)
	}
}

// ResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// StatusCode returns the captured status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}
