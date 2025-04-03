package middleware

import (
	"context"
	"goapi-starter/internal/cache"
	"goapi-starter/internal/config"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"goapi-starter/internal/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Msg("Processing authentication")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Msg("Missing authorization header")
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			logger.Warn().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Str("auth_header", authHeader).
				Msg("Invalid token format")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		tokenStr := bearerToken[1]

		// Check if token is blacklisted - CRITICAL CHECK
		blacklisted, err := cache.IsAccessTokenBlacklisted(tokenStr)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("Error checking token blacklist")
			// If we can't check the blacklist, fail closed for security
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication error")
			return
		}

		if blacklisted {
			logger.Warn().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Msg("Token is blacklisted")
			metrics.RecordHandlerError("AuthMiddleware", "blacklisted_token")
			utils.RespondWithError(w, http.StatusUnauthorized, "Token has been revoked or expired. Please sign in again.")
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Warn().
				Err(err).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Msg("Invalid token")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Msg("Invalid token claims")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Extract user ID from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.Warn().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Msg("Invalid user ID in token")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Create a context with the user ID
		ctx := context.WithValue(r.Context(), "userID", userID)

		// Store the token in context for potential blacklisting during logout
		ctx = context.WithValue(ctx, "accessToken", tokenStr)

		// Try to get user from cache
		cachedUser, found, err := cache.GetCachedUser(userID)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("user_id", userID).
				Msg("Error retrieving user from cache in middleware")
			// Continue with just the user ID in context
		} else if found && cachedUser != nil {
			// Add the full user object to the context if found in cache
			logger.Debug().
				Str("user_id", userID).
				Str("username", cachedUser.Username).
				Msg("User data retrieved from cache in middleware")
			ctx = context.WithValue(ctx, "user", cachedUser)
		}

		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Str("user_id", userID).
			Bool("user_cached", found && cachedUser != nil).
			Msg("Authentication successful")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
