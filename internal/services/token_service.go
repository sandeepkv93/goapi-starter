package services

import (
	"cursor-experiment-1/internal/config"
	"cursor-experiment-1/internal/database"
	"cursor-experiment-1/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenPair(user models.User) (*models.TokenResponse, error) {
	// Generate access token
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.AppConfig.JWT.AccessExpiry,
	}, nil
}

func generateAccessToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Second * time.Duration(config.AppConfig.JWT.AccessExpiry)).Unix(),
		"type":     "access",
	})

	return token.SignedString([]byte(config.AppConfig.JWT.AccessSecret))
}

func generateRefreshToken(user models.User) (string, error) {
	// Generate refresh token string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Second * time.Duration(config.AppConfig.JWT.RefreshExpiry)).Unix(),
		"type":    "refresh",
	})

	refreshTokenString, err := token.SignedString([]byte(config.AppConfig.JWT.RefreshSecret))
	if err != nil {
		return "", err
	}

	// Store refresh token in database
	refreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(time.Second * time.Duration(config.AppConfig.JWT.RefreshExpiry)),
	}

	if result := database.DB.Create(&refreshToken); result.Error != nil {
		return "", result.Error
	}

	return refreshTokenString, nil
}

func ValidateRefreshToken(tokenString string) (*models.User, error) {
	// Find token in database
	var refreshToken models.RefreshToken
	if result := database.DB.Where("token = ? AND expires_at > ?", tokenString, time.Now()).First(&refreshToken); result.Error != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.RefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	var user models.User
	if result := database.DB.First(&user, "id = ?", refreshToken.UserID); result.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
