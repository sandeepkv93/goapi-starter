package utils

import "net/http"

// GetClientIP extracts the client IP address from the request
func GetClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (for proxies)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		for i := 0; i < len(ip); i++ {
			if ip[i] == ',' {
				ip = ip[:i]
				break
			}
		}
		return ip
	}

	// Check for X-Real-IP header (used by some proxies)
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	ip = r.RemoteAddr
	// Remove port if present
	for i := 0; i < len(ip); i++ {
		if ip[i] == ':' {
			ip = ip[:i]
			break
		}
	}

	return ip
}
