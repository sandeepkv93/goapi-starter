package utils

import (
	"context"
	"goapi-starter/internal/models"
)

// GetUserFromContext retrieves the user from the context if available
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value("user").(*models.User)
	return user, ok
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

// GetAccessTokenFromContext retrieves the access token from the context
func GetAccessTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value("accessToken").(string)
	return token, ok
}
