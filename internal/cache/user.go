package cache

import (
	"fmt"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/models"
	"time"
)

const (
	// UserCachePrefix is the prefix for user cache keys
	UserCachePrefix = "user"
	// UserCacheTTL is how long to cache user data (shorter than token expiry)
	UserCacheTTL = 15 * time.Minute
)

// CacheUser stores a user in the cache
func CacheUser(user models.User) error {
	// Don't cache the password
	userToCache := user
	userToCache.Password = ""

	key := fmt.Sprintf("%s:%s", UserCachePrefix, user.ID)
	return SetWithTTL(key, userToCache, UserCacheTTL)
}

// GetCachedUser retrieves a user from the cache
func GetCachedUser(userID string) (*models.User, bool, error) {
	key := fmt.Sprintf("%s:%s", UserCachePrefix, userID)
	var user models.User

	found, err := Get(key, &user)
	if err != nil {
		logger.Warn().Err(err).Str("user_id", userID).Msg("Error retrieving user from cache")
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	return &user, true, nil
}

// InvalidateUserCache removes a user from the cache
func InvalidateUserCache(userID string) error {
	key := fmt.Sprintf("%s:%s", UserCachePrefix, userID)
	return Delete(key)
}
