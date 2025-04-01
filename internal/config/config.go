package config

import (
	"goapi-starter/internal/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  int
	RefreshExpiry int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var AppConfig Config

func LoadConfig() {
	logger.Debug().Msg("Loading application configuration")

	if err := godotenv.Load(); err != nil {
		logger.Warn().Err(err).Msg(".env file not found, using environment variables or defaults")
	} else {
		logger.Debug().Msg("Loaded configuration from .env file")
	}

	AppConfig = Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "3000"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", "default-access-secret"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "default-refresh-secret"),
			AccessExpiry:  getEnvAsInt("JWT_ACCESS_EXPIRY", 900),
			RefreshExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY", 604800),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "goapi_starter_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: loadRedisConfig(),
	}

	// Log configuration (excluding sensitive data)
	logger.Info().
		Str("server_port", AppConfig.Server.Port).
		Int("jwt_access_expiry", AppConfig.JWT.AccessExpiry).
		Int("jwt_refresh_expiry", AppConfig.JWT.RefreshExpiry).
		Str("db_host", AppConfig.Database.Host).
		Str("db_port", AppConfig.Database.Port).
		Str("db_name", AppConfig.Database.DBName).
		Str("db_user", AppConfig.Database.User).
		Str("db_sslmode", AppConfig.Database.SSLMode).
		Msg("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		logger.Debug().Str("key", key).Str("value", value).Msg("Using environment variable")
		return value
	}
	logger.Debug().Str("key", key).Str("default", defaultValue).Msg("Using default value")
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			logger.Debug().Str("key", key).Int("value", intVal).Msg("Using environment variable as int")
			return intVal
		}
		logger.Warn().Str("key", key).Str("value", value).Msg("Failed to parse environment variable as int, using default")
	}
	logger.Debug().Str("key", key).Int("default", defaultValue).Msg("Using default value")
	return defaultValue
}
