package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test default values
	t.Run("Default Values", func(t *testing.T) {
		// Clear any existing env vars that might interfere
		os.Clearenv()

		LoadConfig()

		// Check Server defaults
		if AppConfig.Server.Port != "3000" {
			t.Errorf("Expected default port to be 3000, got %s", AppConfig.Server.Port)
		}

		// Check JWT defaults
		if AppConfig.JWT.AccessSecret != "default-access-secret" {
			t.Errorf("Expected default access secret, got %s", AppConfig.JWT.AccessSecret)
		}
		if AppConfig.JWT.RefreshSecret != "default-refresh-secret" {
			t.Errorf("Expected default refresh secret, got %s", AppConfig.JWT.RefreshSecret)
		}
		if AppConfig.JWT.AccessExpiry != 900 {
			t.Errorf("Expected default access expiry to be 900, got %d", AppConfig.JWT.AccessExpiry)
		}
		if AppConfig.JWT.RefreshExpiry != 604800 {
			t.Errorf("Expected default refresh expiry to be 604800, got %d", AppConfig.JWT.RefreshExpiry)
		}

		// Check Database defaults
		if AppConfig.Database.Host != "localhost" {
			t.Errorf("Expected default host to be localhost, got %s", AppConfig.Database.Host)
		}
		if AppConfig.Database.Port != "5432" {
			t.Errorf("Expected default port to be 5432, got %s", AppConfig.Database.Port)
		}
	})

	// Test custom environment values
	t.Run("Custom Environment Values", func(t *testing.T) {
		os.Clearenv()

		// Set custom environment variables
		os.Setenv("SERVER_PORT", "8080")
		os.Setenv("JWT_ACCESS_SECRET", "custom-access-secret")
		os.Setenv("JWT_ACCESS_EXPIRY", "1800")
		os.Setenv("DB_HOST", "custom-host")

		LoadConfig()

		if AppConfig.Server.Port != "8080" {
			t.Errorf("Expected port to be 8080, got %s", AppConfig.Server.Port)
		}
		if AppConfig.JWT.AccessSecret != "custom-access-secret" {
			t.Errorf("Expected custom access secret, got %s", AppConfig.JWT.AccessSecret)
		}
		if AppConfig.JWT.AccessExpiry != 1800 {
			t.Errorf("Expected access expiry to be 1800, got %d", AppConfig.JWT.AccessExpiry)
		}
		if AppConfig.Database.Host != "custom-host" {
			t.Errorf("Expected host to be custom-host, got %s", AppConfig.Database.Host)
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("Existing Environment Variable", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("TEST_KEY", "test-value")

		value := getEnv("TEST_KEY", "default")
		if value != "test-value" {
			t.Errorf("Expected test-value, got %s", value)
		}
	})

	t.Run("Non-existing Environment Variable", func(t *testing.T) {
		os.Clearenv()

		value := getEnv("NON_EXISTING_KEY", "default")
		if value != "default" {
			t.Errorf("Expected default, got %s", value)
		}
	})
}

func TestGetEnvAsInt(t *testing.T) {
	t.Run("Valid Integer Environment Variable", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("TEST_INT", "123")

		value := getEnvAsInt("TEST_INT", 456)
		if value != 123 {
			t.Errorf("Expected 123, got %d", value)
		}
	})

	t.Run("Invalid Integer Environment Variable", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("TEST_INT", "not-an-int")

		value := getEnvAsInt("TEST_INT", 456)
		if value != 456 {
			t.Errorf("Expected 456, got %d", value)
		}
	})

	t.Run("Non-existing Environment Variable", func(t *testing.T) {
		os.Clearenv()

		value := getEnvAsInt("NON_EXISTING_INT", 456)
		if value != 456 {
			t.Errorf("Expected 456, got %d", value)
		}
	})
}
