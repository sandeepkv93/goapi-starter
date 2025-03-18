package config

const (
	AccessTokenSecret  = "your-access-token-secret"  // Use env vars in production
	RefreshTokenSecret = "your-refresh-token-secret" // Use env vars in production
	AccessTokenExpiry  = 15 * 60                     // 15 minutes in seconds
	RefreshTokenExpiry = 7 * 24 * 60 * 60            // 7 days in seconds
)
