package handlers

import (
	"goapi-starter/internal/cache"
	"goapi-starter/internal/database"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"goapi-starter/internal/models"
	"goapi-starter/internal/utils"
	"net/http"
)

// GetProfile returns the current user's profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("get_profile", "started").Inc()

	// Try to get user from context (cached)
	user, found := utils.GetUserFromContext(r.Context())

	if found {
		logger.Debug().
			Str("user_id", user.ID).
			Str("username", user.Username).
			Msg("User retrieved from context cache")

		metrics.BusinessOperations.WithLabelValues("get_profile", "success").Inc()
		utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
			Message: "Profile retrieved from cache",
			Data:    user,
		})
		return
	}

	// If not found in context, get user ID and fetch from database
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		metrics.RecordHandlerError("GetProfile", "unauthorized")
		metrics.BusinessOperations.WithLabelValues("get_profile", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Try to get from cache directly (fallback)
	cachedUser, found, err := cache.GetCachedUser(userID)
	if err == nil && found && cachedUser != nil {
		logger.Debug().
			Str("user_id", cachedUser.ID).
			Str("username", cachedUser.Username).
			Msg("User retrieved from Redis cache")

		metrics.BusinessOperations.WithLabelValues("get_profile", "success").Inc()
		utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
			Message: "Profile retrieved from Redis cache",
			Data:    cachedUser,
		})
		return
	}

	// Not in cache, get from database
	var dbUser models.User
	if result := database.DB.First(&dbUser, "id = ?", userID); result.Error != nil {
		metrics.RecordHandlerError("GetProfile", "user_not_found")
		metrics.BusinessOperations.WithLabelValues("get_profile", "failed").Inc()
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Cache the user for future requests
	if err := cache.CacheUser(dbUser); err != nil {
		logger.Warn().
			Err(err).
			Str("user_id", dbUser.ID).
			Msg("Failed to cache user data")
		// Continue even if caching fails
	}

	// Don't return the password
	dbUser.Password = ""

	metrics.BusinessOperations.WithLabelValues("get_profile", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Profile retrieved from database",
		Data:    dbUser,
	})
}
