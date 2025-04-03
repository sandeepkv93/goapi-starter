package cache

import (
	"fmt"
	"goapi-starter/internal/config"
	"goapi-starter/internal/logger"
	"time"
)

const (
	// TokenBlacklistPrefix is the prefix for blacklisted tokens
	TokenBlacklistPrefix = "blacklist"
)

// BlacklistToken adds a token to the blacklist
// The TTL should match the token's remaining validity period
func BlacklistToken(tokenType string, token string, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s:%s", TokenBlacklistPrefix, tokenType, token)

	// Store a simple value (timestamp when it was blacklisted)
	now := time.Now().Unix()

	logger.Debug().
		Str("token_type", tokenType).
		Str("token", token[:10]+"..."). // Only log part of the token for security
		Dur("ttl", ttl).
		Msg("Blacklisting token")

	return SetWithTTL(key, now, ttl)
}

// IsTokenBlacklisted checks if a token is in the blacklist
func IsTokenBlacklisted(tokenType string, token string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", TokenBlacklistPrefix, tokenType, token)

	var timestamp int64
	found, err := Get(key, &timestamp)

	if err != nil {
		logger.Warn().
			Err(err).
			Str("token_type", tokenType).
			Str("token", token[:10]+"...").
			Msg("Error checking token blacklist")
		return false, err
	}

	return found, nil
}

// BlacklistAccessToken adds an access token to the blacklist
func BlacklistAccessToken(token string) error {
	// Use the configured access token expiry or a default
	ttl := time.Duration(config.AppConfig.JWT.AccessExpiry) * time.Second
	if ttl <= 0 {
		ttl = 15 * time.Minute // Default if not configured
	}

	return BlacklistToken("access", token, ttl)
}

// BlacklistRefreshToken adds a refresh token to the blacklist
func BlacklistRefreshToken(token string) error {
	// Use the configured refresh token expiry or a default
	ttl := time.Duration(config.AppConfig.JWT.RefreshExpiry) * time.Second
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour // Default if not configured
	}

	return BlacklistToken("refresh", token, ttl)
}

// IsAccessTokenBlacklisted checks if an access token is blacklisted
func IsAccessTokenBlacklisted(token string) (bool, error) {
	return IsTokenBlacklisted("access", token)
}

// IsRefreshTokenBlacklisted checks if a refresh token is blacklisted
func IsRefreshTokenBlacklisted(token string) (bool, error) {
	return IsTokenBlacklisted("refresh", token)
}
