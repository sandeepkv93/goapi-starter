package middleware

import (
	"context"
	"goapi-starter/internal/config"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn().Msg("Missing authorization header")
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			logger.Warn().Msg("Invalid token format")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		tokenStr := bearerToken[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Warn().Err(err).Msg("Invalid token")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn().Msg("Invalid token claims")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "userID", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
