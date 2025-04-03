package utils

import (
	"encoding/json"
	"goapi-starter/internal/logger"
	"net/http"
)

type ErrorResponse struct {
	Error         string `json:"error"`
	CorrelationID string `json:"correlation_id"`
}

type SuccessResponse struct {
	Message       string      `json:"message"`
	Data          interface{} `json:"data,omitempty"`
	CorrelationID string      `json:"correlation_id"`
}

func RespondWithError(w http.ResponseWriter, r *http.Request, code int, message string) {
	logger.Debug().
		Int("status_code", code).
		Str("error", message).
		Msg("Sending error response")

	// Get correlation ID from request context
	correlationID := GetCorrelationID(r.Context())

	RespondWithJSON(w, r, code, ErrorResponse{
		Error:         message,
		CorrelationID: correlationID,
	})
}

func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	// If payload is a SuccessResponse or ErrorResponse, add correlation ID
	switch v := payload.(type) {
	case SuccessResponse:
		if v.CorrelationID == "" {
			v.CorrelationID = GetCorrelationID(r.Context())
			payload = v
		}
	case ErrorResponse:
		if v.CorrelationID == "" {
			v.CorrelationID = GetCorrelationID(r.Context())
			payload = v
		}
	}

	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error().
			Err(err).
			Int("status_code", code).
			Interface("payload", payload).
			Msg("Error marshalling JSON response")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling JSON"}`))
		return
	}

	logger.Debug().
		Int("status_code", code).
		Int("response_size", len(response)).
		Msg("Sending JSON response")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
