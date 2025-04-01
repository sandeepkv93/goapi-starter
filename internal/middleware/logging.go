package middleware

import (
	"goapi-starter/internal/logger"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status       int
	wroteHeader  bool
	bytesWritten int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
		rw.wroteHeader = true
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func (rw *responseWriter) BytesWritten() int {
	return rw.bytesWritten
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		// Log request details before processing
		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("Request received")

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Choose log level based on status code
		logEvent := logger.Info()
		if wrapped.status >= 400 && wrapped.status < 500 {
			logEvent = logger.Warn()
		} else if wrapped.status >= 500 {
			logEvent = logger.Error()
		}

		// Log the request details
		logEvent.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Int("status", wrapped.status).
			Int("bytes", wrapped.BytesWritten()).
			Dur("duration", duration).
			Str("duration_human", duration.String()).
			Msg("Request completed")
	})
}
