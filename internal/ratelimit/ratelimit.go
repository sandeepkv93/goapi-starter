package ratelimit

import (
	"fmt"
	"goapi-starter/internal/cache"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"time"
)

const (
	// Default limits
	DefaultIPRateLimit   = 60  // requests per minute for unauthenticated users
	DefaultUserRateLimit = 300 // requests per minute for authenticated users
	DefaultAuthRateLimit = 10  // login/signup attempts per minute
	DefaultWindowSize    = 60  // seconds (1 minute window)
	DefaultBlockDuration = 300 // seconds (5 minute block after exceeding limit)

	// Redis key prefixes
	IPLimitPrefix   = "ratelimit:ip:"
	UserLimitPrefix = "ratelimit:user:"
	AuthLimitPrefix = "ratelimit:auth:"
)

// RateLimiter defines the configuration for rate limiting
type RateLimiter struct {
	// Maximum number of requests allowed in the time window
	Limit int
	// Time window in seconds
	WindowSize int
	// How long to block after exceeding limit (seconds)
	BlockDuration int
	// Key prefix for Redis
	KeyPrefix string
}

// NewIPRateLimiter creates a rate limiter for IP-based limiting
func NewIPRateLimiter() *RateLimiter {
	return &RateLimiter{
		Limit:         DefaultIPRateLimit,
		WindowSize:    DefaultWindowSize,
		BlockDuration: DefaultBlockDuration,
		KeyPrefix:     IPLimitPrefix,
	}
}

// NewUserRateLimiter creates a rate limiter for user-based limiting
func NewUserRateLimiter() *RateLimiter {
	return &RateLimiter{
		Limit:         DefaultUserRateLimit,
		WindowSize:    DefaultWindowSize,
		BlockDuration: DefaultBlockDuration,
		KeyPrefix:     UserLimitPrefix,
	}
}

// NewAuthRateLimiter creates a stricter rate limiter for authentication endpoints
func NewAuthRateLimiter() *RateLimiter {
	return &RateLimiter{
		Limit:         DefaultAuthRateLimit,
		WindowSize:    DefaultWindowSize,
		BlockDuration: DefaultBlockDuration,
		KeyPrefix:     AuthLimitPrefix,
	}
}

// Allow checks if a request should be allowed based on the rate limit
// Returns: allowed (bool), remaining (int), resetAfter (time.Duration), err (error)
func (rl *RateLimiter) Allow(identifier string) (bool, int, time.Duration, error) {
	metrics.RecordRateLimitCheck(rl.KeyPrefix)

	// Check if the identifier is currently blocked
	blockedKey := fmt.Sprintf("%s%s:blocked", rl.KeyPrefix, identifier)
	var blockedUntil int64
	blocked, err := cache.Get(blockedKey, &blockedUntil)

	if err != nil {
		logger.Warn().
			Err(err).
			Str("identifier", identifier).
			Msg("Error checking rate limit block status")
		// On error, we'll continue and assume not blocked
	} else if blocked {
		now := time.Now().Unix()
		if blockedUntil > now {
			// Still blocked
			resetAfter := time.Duration(blockedUntil-now) * time.Second
			metrics.RecordRateLimitResult(rl.KeyPrefix, "blocked")
			return false, 0, resetAfter, nil
		}
		// Block has expired, remove it
		_ = cache.Delete(blockedKey)
	}

	// Get the current count
	key := fmt.Sprintf("%s%s", rl.KeyPrefix, identifier)
	var count int
	exists, err := cache.Get(key, &count)

	if err != nil {
		logger.Warn().
			Err(err).
			Str("identifier", identifier).
			Msg("Error checking rate limit count")
		// On error, we'll allow the request but not increment
		metrics.RecordRateLimitResult(rl.KeyPrefix, "error")
		return true, rl.Limit, time.Duration(rl.WindowSize) * time.Second, err
	}

	now := time.Now()
	windowExpiry := time.Duration(rl.WindowSize) * time.Second

	if !exists {
		// First request in this window
		count = 1
		err = cache.SetWithTTL(key, count, windowExpiry)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("identifier", identifier).
				Msg("Error setting initial rate limit count")
		}
		metrics.RecordRateLimitResult(rl.KeyPrefix, "allowed")
		return true, rl.Limit - 1, windowExpiry, nil
	}

	// Increment the counter
	count++

	// Check if limit exceeded
	if count > rl.Limit {
		// Block the identifier for the block duration
		blockUntil := now.Add(time.Duration(rl.BlockDuration) * time.Second).Unix()
		err = cache.SetWithTTL(blockedKey, blockUntil, time.Duration(rl.BlockDuration)*time.Second)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("identifier", identifier).
				Msg("Error setting rate limit block")
		}

		metrics.RecordRateLimitResult(rl.KeyPrefix, "exceeded")
		return false, 0, time.Duration(rl.BlockDuration) * time.Second, nil
	}

	// Update the counter
	// Get TTL of the existing key to maintain the same window
	ttl, err := cache.GetTTL(key)
	if err != nil || ttl <= 0 {
		ttl = windowExpiry
	}

	err = cache.SetWithTTL(key, count, ttl)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("identifier", identifier).
			Int("count", count).
			Msg("Error updating rate limit count")
	}

	metrics.RecordRateLimitResult(rl.KeyPrefix, "allowed")
	return true, rl.Limit - count, ttl, nil
}
