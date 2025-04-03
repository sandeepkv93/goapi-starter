package handlers

import (
	"encoding/json"
	"goapi-starter/internal/cache"
	"goapi-starter/internal/database"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"goapi-starter/internal/models"
	"goapi-starter/internal/services"
	"goapi-starter/internal/utils"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("signup", "started").Inc()

	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("SignUp", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("SignUp", "validation_error")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if email already exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		metrics.RecordHandlerError("SignUp", "email_exists")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	// Check if username already exists
	if result := database.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		metrics.RecordHandlerError("SignUp", "username_exists")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusConflict, "Username already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		metrics.RecordHandlerError("SignUp", "password_hash_error")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error processing request")
		return
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if result := database.DB.Create(&user); result.Error != nil {
		metrics.RecordHandlerError("SignUp", "database_error")
		metrics.BusinessOperations.WithLabelValues("signup", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	metrics.BusinessOperations.WithLabelValues("signup", "success").Inc()
	// Return user data without password
	utils.RespondWithJSON(w, http.StatusCreated, utils.SuccessResponse{
		Message: "User created successfully",
		Data: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("signin", "started").Inc()

	var req models.SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("SignIn", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("SignIn", "validation_error")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Find user by email
	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		metrics.RecordHandlerError("SignIn", "user_not_found")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		metrics.RecordHandlerError("SignIn", "invalid_password")
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token pair
	tokens, err := services.GenerateTokenPair(user)
	if err != nil {
		errorReason := "unknown"
		if err.Error() != "" {
			// Extract a simplified error reason
			if strings.Contains(err.Error(), "duplicate key value") {
				errorReason = "duplicate_token"
			} else if strings.Contains(err.Error(), "database") {
				errorReason = "database_error"
			} else {
				// Limit the error reason length to avoid cardinality explosion
				if len(err.Error()) > 50 {
					errorReason = err.Error()[:50]
				} else {
					errorReason = err.Error()
				}
			}
		}

		metrics.RecordHandlerError("SignIn", "token_generation_error")
		metrics.RecordDetailedError("SignIn", "token_generation_error", errorReason)
		metrics.BusinessOperations.WithLabelValues("signin", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating tokens")
		return
	}

	// Cache the user for future requests
	if err := cache.CacheUser(user); err != nil {
		logger.Warn().Err(err).Str("user_id", user.ID).Msg("Failed to cache user data")
		// Continue even if caching fails
	} else {
		logger.Debug().Str("user_id", user.ID).Msg("User data cached successfully")
	}

	metrics.BusinessOperations.WithLabelValues("signin", "success").Inc()

	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Successfully signed in",
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
			},
			"tokens": tokens,
		},
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("refresh_token", "started").Inc()

	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordHandlerError("RefreshToken", "invalid_request")
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		metrics.RecordHandlerError("RefreshToken", "validation_error")
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate refresh token and get user
	user, err := services.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		errorReason := "unknown"
		if err.Error() != "" {
			if strings.Contains(err.Error(), "expired") {
				errorReason = "token_expired"
			} else if strings.Contains(err.Error(), "not found") {
				errorReason = "token_not_found"
			} else if strings.Contains(err.Error(), "invalid") {
				errorReason = "token_invalid"
			} else {
				// Limit the error reason length to avoid cardinality explosion
				if len(err.Error()) > 50 {
					errorReason = err.Error()[:50]
				} else {
					errorReason = err.Error()
				}
			}
		}

		metrics.RecordHandlerError("RefreshToken", "invalid_token")
		metrics.RecordDetailedError("RefreshToken", "invalid_token", errorReason)
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Generate new token pair
	tokens, err := services.GenerateTokenPair(*user)
	if err != nil {
		errorReason := "unknown"
		if err.Error() != "" {
			if strings.Contains(err.Error(), "duplicate key value") {
				errorReason = "duplicate_token"
			} else if strings.Contains(err.Error(), "database") {
				errorReason = "database_error"
			} else {
				// Limit the error reason length to avoid cardinality explosion
				if len(err.Error()) > 50 {
					errorReason = err.Error()[:50]
				} else {
					errorReason = err.Error()
				}
			}
		}

		metrics.RecordHandlerError("RefreshToken", "token_generation_error")
		metrics.RecordDetailedError("RefreshToken", "token_generation_error", errorReason)
		metrics.BusinessOperations.WithLabelValues("refresh_token", "failed").Inc()
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating tokens")
		return
	}

	metrics.BusinessOperations.WithLabelValues("refresh_token", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Tokens refreshed successfully",
		Data:    tokens,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	metrics.BusinessOperations.WithLabelValues("logout", "started").Inc()

	// Get user ID from context
	userID, userIDFound := utils.GetUserIDFromContext(r.Context())
	if !userIDFound || userID == "" {
		metrics.RecordHandlerError("Logout", "user_not_found")
		metrics.BusinessOperations.WithLabelValues("logout", "failed").Inc()
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get the current access token from context
	accessToken, ok := r.Context().Value("accessToken").(string)
	if ok && accessToken != "" {
		logger.Debug().
			Str("token_prefix", accessToken[:10]+"...").
			Msg("Blacklisting access token from context")

		// Blacklist the access token
		if err := cache.BlacklistAccessToken(accessToken); err != nil {
			logger.Warn().
				Err(err).
				Msg("Failed to blacklist access token")
			// Continue even if blacklisting fails
		} else {
			logger.Info().
				Str("token_prefix", accessToken[:10]+"...").
				Msg("Access token blacklisted successfully")
		}
	} else {
		logger.Warn().Msg("No access token found in context during logout")
	}

	// Invalidate user cache
	if err := cache.InvalidateUserCache(userID); err != nil {
		logger.Warn().
			Err(err).
			Str("user_id", userID).
			Msg("Failed to invalidate user cache during logout")
		// Continue even if cache invalidation fails
	} else {
		logger.Debug().
			Str("user_id", userID).
			Msg("User cache invalidated during logout")
	}

	// Find and invalidate all active refresh tokens for this user
	var refreshTokens []models.RefreshToken
	if result := database.DB.Where("user_id = ? AND expires_at > NOW()", userID).Find(&refreshTokens); result.Error != nil {
		logger.Warn().
			Err(result.Error).
			Str("user_id", userID).
			Msg("Failed to retrieve refresh tokens during logout")
	} else {
		// Blacklist all active refresh tokens
		for _, rt := range refreshTokens {
			if err := cache.BlacklistRefreshToken(rt.Token); err != nil {
				logger.Warn().
					Err(err).
					Str("token_id", rt.ID).
					Msg("Failed to blacklist refresh token")
			} else {
				logger.Debug().
					Str("token_id", rt.ID).
					Msg("Refresh token blacklisted successfully")
			}
		}

		// Delete all refresh tokens for this user
		if result := database.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}); result.Error != nil {
			logger.Warn().
				Err(result.Error).
				Str("user_id", userID).
				Msg("Failed to delete refresh tokens during logout")
		} else {
			logger.Debug().
				Int64("count", result.RowsAffected).
				Str("user_id", userID).
				Msg("Refresh tokens deleted successfully")
		}
	}

	metrics.BusinessOperations.WithLabelValues("logout", "success").Inc()
	utils.RespondWithJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Logged out successfully",
	})
}
