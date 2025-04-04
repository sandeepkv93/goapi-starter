package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"goapi-starter/internal/config"
	"goapi-starter/internal/database"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/models"
	"time"

	"goapi-starter/internal/cache"

	"github.com/golang-jwt/jwt/v5"
)

// generateRandomString creates a random string for token uniqueness
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		logger.Error().Err(err).Int("length", length).Msg("Failed to generate random string")
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateTokenPair(user models.User) (*models.TokenResponse, error) {
	logger.Debug().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Msg("Generating token pair")

	// Generate access token
	accessToken, err := generateAccessToken(user)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID).
			Str("username", user.Username).
			Msg("Failed to generate access token")
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID).
			Str("username", user.Username).
			Msg("Failed to generate refresh token")
		return nil, err
	}

	logger.Info().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Int("expires_in", config.AppConfig.JWT.AccessExpiry).
		Msg("Token pair generated successfully")

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.AppConfig.JWT.AccessExpiry,
	}, nil
}

func generateAccessToken(user models.User) (string, error) {
	logger.Debug().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Int("expiry", config.AppConfig.JWT.AccessExpiry).
		Msg("Generating access token")

	expiryTime := time.Now().Add(time.Second * time.Duration(config.AppConfig.JWT.AccessExpiry))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expiryTime.Unix(),
		"type":     "access",
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.AccessSecret))
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("Failed to sign access token")
		return "", err
	}

	logger.Debug().
		Str("user_id", user.ID).
		Time("expires_at", expiryTime).
		Msg("Access token generated successfully")

	return tokenString, nil
}

func generateRefreshToken(user models.User) (string, error) {
	logger.Debug().
		Str("user_id", user.ID).
		Int("expiry", config.AppConfig.JWT.RefreshExpiry).
		Msg("Generating refresh token")

	// Generate a random component to ensure uniqueness
	randomID, err := generateRandomString(16)
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("Failed to generate random string for refresh token")
		return "", err
	}

	expiryTime := time.Now().Add(time.Second * time.Duration(config.AppConfig.JWT.RefreshExpiry))

	// Generate refresh token string with the random component
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"jti":     randomID, // Add a unique JWT ID
		"exp":     expiryTime.Unix(),
		"type":    "refresh",
	})

	refreshTokenString, err := token.SignedString([]byte(config.AppConfig.JWT.RefreshSecret))
	if err != nil {
		logger.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("Failed to sign refresh token")
		return "", err
	}

	// Store refresh token in database
	refreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiryTime,
	}

	if result := database.DB.Create(&refreshToken); result.Error != nil {
		logger.Error().
			Err(result.Error).
			Str("user_id", user.ID).
			Time("expires_at", expiryTime).
			Msg("Failed to store refresh token in database")
		return "", result.Error
	}

	logger.Debug().
		Str("user_id", user.ID).
		Time("expires_at", expiryTime).
		Msg("Refresh token generated and stored successfully")

	return refreshTokenString, nil
}

func ValidateRefreshToken(tokenString string) (*models.User, error) {
	logger.Debug().Msg("Validating refresh token")

	// Check if token is blacklisted
	blacklisted, err := cache.IsRefreshTokenBlacklisted(tokenString)
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Error checking refresh token blacklist")
		// Continue even if blacklist check fails
	} else if blacklisted {
		logger.Warn().
			Msg("Refresh token is blacklisted")
		return nil, errors.New("token has been revoked")
	}

	// Check if token is in Redis cache
	userID, found, err := cache.GetCachedRefreshToken(tokenString)
	if err == nil && found && userID != "" {
		logger.Debug().
			Str("user_id", userID).
			Msg("Refresh token found in cache, skipping database validation")

		// Token is valid and cached, get user
		return getUserForToken(userID)
	}

	// Not in cache, validate from database
	var refreshToken models.RefreshToken
	if result := database.DB.Where("token = ? AND expires_at > ?", tokenString, time.Now()).First(&refreshToken); result.Error != nil {
		logger.Warn().
			Err(result.Error).
			Msg("Refresh token not found or expired in database")
		return nil, errors.New("invalid refresh token")
	}

	logger.Debug().
		Str("user_id", refreshToken.UserID).
		Time("expires_at", refreshToken.ExpiresAt).
		Msg("Refresh token found in database, validating JWT")

	// Validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.RefreshSecret), nil
	})

	if err != nil || !token.Valid {
		logger.Warn().
			Err(err).
			Str("user_id", refreshToken.UserID).
			Msg("Invalid JWT refresh token")
		return nil, errors.New("invalid refresh token")
	}

	// Cache the validated token for future checks
	timeUntilExpiry := time.Until(refreshToken.ExpiresAt)
	if timeUntilExpiry > 0 {
		if cacheErr := cache.CacheRefreshToken(tokenString, refreshToken.UserID, timeUntilExpiry); cacheErr != nil {
			logger.Warn().
				Err(cacheErr).
				Str("user_id", refreshToken.UserID).
				Msg("Failed to cache validated refresh token")
			// Continue even if caching fails
		}
	}

	// Get user for the token
	return getUserForToken(refreshToken.UserID)
}

// Helper function to get user by ID (from cache or database)
func getUserForToken(userID string) (*models.User, error) {
	// Try to get user from cache first
	cachedUser, found, err := cache.GetCachedUser(userID)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("user_id", userID).
			Msg("Error retrieving user from cache")
		// Continue with database lookup
	}

	if found && cachedUser != nil {
		logger.Info().
			Str("user_id", cachedUser.ID).
			Str("username", cachedUser.Username).
			Msg("User retrieved from cache for refresh token")
		return cachedUser, nil
	}

	// User not in cache, get from database
	var user models.User
	if result := database.DB.First(&user, "id = ?", userID); result.Error != nil {
		logger.Error().
			Err(result.Error).
			Str("user_id", userID).
			Msg("User not found for refresh token")
		return nil, errors.New("user not found")
	}

	// Cache the user for future requests
	if err := cache.CacheUser(user); err != nil {
		logger.Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("Failed to cache user data during token validation")
		// Continue even if caching fails
	}

	logger.Info().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Msg("Refresh token validated successfully")

	return &user, nil
}
