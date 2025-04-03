package middleware

import (
	"goapi-starter/internal/logger"
	"goapi-starter/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func CorrelationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate or get correlation ID
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Get user ID if available
		userID, _ := utils.GetUserIDFromContext(r.Context())

		// Create a new logger with correlation context
		contextLogger := log.Logger.With().
			Str("correlation_id", correlationID)

		if userID != "" {
			contextLogger = contextLogger.Str("user_id", userID)
		}

		// Replace the default logger for this request
		logger.SetRequestLogger(contextLogger.Logger())

		defer logger.ClearRequestLogger() // Clean up after request is done

		next.ServeHTTP(w, r)
	})
}
