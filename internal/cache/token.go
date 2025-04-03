package cache

import (
	"fmt"
	"goapi-starter/internal/logger"
	"time"
)

const (
	// RefreshTokenCachePrefix is the prefix for cached refresh tokens
	RefreshTokenCachePrefix = "refresh_token"
	// RefreshTokenCacheTTL is how long to cache validated refresh tokens
	RefreshTokenCacheTTL = 24 * time.Hour
)

// CacheRefreshToken stores a validated refresh token in the cache
func CacheRefreshToken(tokenString string, userID string, expiry time.Duration) error {
	// If the provided expiry is longer than our max cache TTL, use the max TTL
	if expiry > RefreshTokenCacheTTL {
		expiry = RefreshTokenCacheTTL
	}

	key := fmt.Sprintf("%s:%s", RefreshTokenCachePrefix, tokenString)
	err := SetWithTTL(key, userID, expiry)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("token_prefix", tokenString[:10]+"...").
			Msg("Failed to cache refresh token")
		return err
	}

	logger.Debug().
		Str("token_prefix", tokenString[:10]+"...").
		Str("user_id", userID).
		Dur("ttl", expiry).
		Msg("Refresh token cached successfully")
	return nil
}

// GetCachedRefreshToken retrieves a refresh token from the cache
func GetCachedRefreshToken(tokenString string) (string, bool, error) {
	key := fmt.Sprintf("%s:%s", RefreshTokenCachePrefix, tokenString)
	var userID string

	found, err := Get(key, &userID)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("token_prefix", tokenString[:10]+"...").
			Msg("Error retrieving refresh token from cache")
		return "", false, err
	}

	if !found || userID == "" {
		return "", false, nil
	}

	logger.Debug().
		Str("token_prefix", tokenString[:10]+"...").
		Str("user_id", userID).
		Msg("Refresh token found in cache")
	return userID, true, nil
}

// InvalidateRefreshTokenCache removes a refresh token from the cache
func InvalidateRefreshTokenCache(tokenString string) error {
	key := fmt.Sprintf("%s:%s", RefreshTokenCachePrefix, tokenString)
	return Delete(key)
}
