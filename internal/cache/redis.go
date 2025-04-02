package cache

import (
	"context"
	"encoding/json"
	"goapi-starter/internal/config"
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis initializes the Redis client
func InitRedis() {
	logger.Info().Msg("Initializing Redis connection")

	redisConfig := config.AppConfig.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + redisConfig.Port,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	logger.Info().
		Str("host", redisConfig.Host).
		Str("port", redisConfig.Port).
		Msg("Successfully connected to Redis")
}

// Set stores a value in the cache with the default TTL
func Set(key string, value interface{}) error {
	metrics.RecordCacheOperation("set", "default")
	return SetWithTTL(key, value, config.AppConfig.Redis.CacheTTL)
}

// SetWithTTL stores a value in the cache with a specific TTL
func SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	metrics.RecordCacheOperation("set", "custom_ttl")
	startTime := time.Now()
	defer func() {
		metrics.RecordCacheDuration("set", time.Since(startTime))
	}()

	// Marshal the value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		logger.Error().Err(err).Str("key", key).Msg("Failed to marshal value for caching")
		return err
	}

	// Record the size of the cached object
	keyPrefix := strings.Split(key, ":")[0]
	metrics.RecordCacheSize(keyPrefix, len(jsonValue))

	// Store in Redis
	err = RedisClient.Set(ctx, key, jsonValue, ttl).Err()
	if err != nil {
		logger.Error().Err(err).Str("key", key).Msg("Failed to set cache value")
		return err
	}

	logger.Debug().Str("key", key).Dur("ttl", ttl).Msg("Successfully cached value")
	return nil
}

// Get retrieves a value from the cache
func Get(key string, dest interface{}) (bool, error) {
	metrics.RecordCacheOperation("get", "default")
	startTime := time.Now()
	defer func() {
		metrics.RecordCacheDuration("get", time.Since(startTime))
	}()

	// Get from Redis
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist
		logger.Debug().Str("key", key).Msg("Cache miss")
		metrics.RecordCacheResult("miss")
		return false, nil
	} else if err != nil {
		// Error occurred
		logger.Error().Err(err).Str("key", key).Msg("Error retrieving from cache")
		metrics.RecordCacheResult("error")
		return false, err
	}

	// Unmarshal the value
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		logger.Error().Err(err).Str("key", key).Msg("Failed to unmarshal cached value")
		metrics.RecordCacheResult("unmarshal_error")
		return false, err
	}

	logger.Debug().Str("key", key).Msg("Cache hit")
	metrics.RecordCacheResult("hit")
	return true, nil
}

// Delete removes a value from the cache
func Delete(key string) error {
	metrics.RecordCacheOperation("delete", "default")
	startTime := time.Now()
	defer func() {
		metrics.RecordCacheDuration("delete", time.Since(startTime))
	}()

	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		logger.Error().Err(err).Str("key", key).Msg("Failed to delete from cache")
		return err
	}

	logger.Debug().Str("key", key).Msg("Successfully deleted from cache")
	return nil
}

// FlushAll clears the entire cache
func FlushAll() error {
	metrics.RecordCacheOperation("flush", "all")
	startTime := time.Now()
	defer func() {
		metrics.RecordCacheDuration("flush", time.Since(startTime))
	}()

	err := RedisClient.FlushAll(ctx).Err()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to flush cache")
		return err
	}

	logger.Info().Msg("Successfully flushed entire cache")
	return nil
}
