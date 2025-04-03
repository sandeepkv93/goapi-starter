package middleware

import (
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"goapi-starter/internal/ratelimit"
	"goapi-starter/internal/utils"
	"net/http"
	"strconv"
)

// IPRateLimitMiddleware limits requests based on IP address
func IPRateLimitMiddleware(next http.Handler) http.Handler {
	limiter := ratelimit.NewIPRateLimiter()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		clientIP := utils.GetClientIP(r)

		allowed, remaining, resetAfter, err := limiter.Allow(clientIP)
		if err != nil {
			// On error, we'll allow the request but log the issue
			logger.Warn().
				Err(err).
				Str("ip", clientIP).
				Str("path", r.URL.Path).
				Msg("Rate limit check error")
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.Limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(int64(resetAfter.Seconds()), 10))

		if !allowed {
			metrics.RecordHandlerError("IPRateLimitMiddleware", "rate_limited")
			logger.Warn().
				Str("ip", clientIP).
				Str("path", r.URL.Path).
				Msg("IP rate limit exceeded")

			w.Header().Set("Retry-After", strconv.FormatInt(int64(resetAfter.Seconds()), 10))
			utils.RespondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UserRateLimitMiddleware limits requests based on user ID for authenticated users
func UserRateLimitMiddleware(next http.Handler) http.Handler {
	limiter := ratelimit.NewUserRateLimiter()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get user ID from context
		userID, ok := utils.GetUserIDFromContext(r.Context())

		// If no user ID (not authenticated), fall back to IP-based limiting
		if !ok || userID == "" {
			// Pass to the next handler - should be caught by auth middleware
			// or handled by the IP rate limiter
			next.ServeHTTP(w, r)
			return
		}

		allowed, remaining, resetAfter, err := limiter.Allow(userID)
		if err != nil {
			// On error, we'll allow the request but log the issue
			logger.Warn().
				Err(err).
				Str("user_id", userID).
				Str("path", r.URL.Path).
				Msg("User rate limit check error")
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.Limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(int64(resetAfter.Seconds()), 10))

		if !allowed {
			metrics.RecordHandlerError("UserRateLimitMiddleware", "rate_limited")
			logger.Warn().
				Str("user_id", userID).
				Str("path", r.URL.Path).
				Msg("User rate limit exceeded")

			w.Header().Set("Retry-After", strconv.FormatInt(int64(resetAfter.Seconds()), 10))
			utils.RespondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthRateLimitMiddleware provides stricter rate limiting for authentication endpoints
func AuthRateLimitMiddleware(next http.Handler) http.Handler {
	limiter := ratelimit.NewAuthRateLimiter()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		clientIP := utils.GetClientIP(r)

		allowed, remaining, resetAfter, err := limiter.Allow(clientIP)
		if err != nil {
			// On error, we'll allow the request but log the issue
			logger.Warn().
				Err(err).
				Str("ip", clientIP).
				Str("path", r.URL.Path).
				Msg("Auth rate limit check error")
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.Limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(int64(resetAfter.Seconds()), 10))

		if !allowed {
			metrics.RecordHandlerError("AuthRateLimitMiddleware", "rate_limited")
			logger.Warn().
				Str("ip", clientIP).
				Str("path", r.URL.Path).
				Msg("Auth rate limit exceeded")

			w.Header().Set("Retry-After", strconv.FormatInt(int64(resetAfter.Seconds()), 10))
			utils.RespondWithError(w, http.StatusTooManyRequests, "Too many authentication attempts. Please try again later.")
			return
		}

		next.ServeHTTP(w, r)
	})
}
