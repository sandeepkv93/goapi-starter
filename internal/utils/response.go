package utils

import (
	"encoding/json"
	"goapi-starter/internal/logger"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	logger.Debug().
		Int("status_code", code).
		Str("error", message).
		Msg("Sending error response")

	RespondWithJSON(w, code, ErrorResponse{Error: message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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
