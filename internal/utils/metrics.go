package utils

import (
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"net/http"
	"time"
)

// InstrumentHandler wraps an HTTP handler with metrics instrumentation
func InstrumentHandler(handlerName string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Debug().
			Str("handler", handlerName).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("Handler execution started")

		// Create a custom response writer to capture the status code
		rw := NewResponseWriter(w)

		// Execute the handler
		next(rw, r)

		// Record metrics
		duration := time.Since(start)
		metrics.RecordHandlerExecution(handlerName, rw.StatusCode(), duration)

		// Log at appropriate level based on status code
		logEvent := logger.Info()
		if rw.statusCode >= 400 && rw.statusCode < 500 {
			logEvent = logger.Warn()
		} else if rw.statusCode >= 500 {
			logEvent = logger.Error()
		}

		logEvent.
			Str("handler", handlerName).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rw.StatusCode()).
			Dur("duration", duration).
			Str("duration_human", duration.String()).
			Msg("Handler execution completed")
	}
}

// ResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK, 0}
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// StatusCode returns the captured status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// BytesWritten returns the number of bytes written
func (rw *ResponseWriter) BytesWritten() int {
	return rw.bytesWritten
}
