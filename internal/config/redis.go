package config

import (
	"goapi-starter/internal/logger"
	"strconv"
	"time"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	CacheTTL time.Duration // Default TTL for cached items
}

func loadRedisConfig() RedisConfig {
	logger.Debug().Msg("Loading Redis configuration")

	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	cacheTTL, _ := strconv.Atoi(getEnv("REDIS_CACHE_TTL", "3600"))

	config := RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
		CacheTTL: time.Duration(cacheTTL) * time.Second,
	}

	logger.Info().
		Str("redis_host", config.Host).
		Str("redis_port", config.Port).
		Int("redis_db", config.DB).
		Dur("cache_ttl", config.CacheTTL).
		Msg("Redis configuration loaded")

	return config
}
