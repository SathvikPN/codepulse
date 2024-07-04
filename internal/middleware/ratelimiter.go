package middleware

import (
	"net/http"
	"sync"
	"time"
)

func RateLimiting(next http.Handler) http.Handler {
	var rateLimit = 2 // Limit to 10 requests per minute
	var interval = time.Minute
	var requestCounts = make(map[string]int)
	var lastReset = time.Now()
	var mu sync.Mutex

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		ip := r.RemoteAddr
		if time.Since(lastReset) > interval {
			lastReset = time.Now()
			requestCounts = make(map[string]int)
		}

		if requestCounts[ip] >= rateLimit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		requestCounts[ip]++
		next.ServeHTTP(w, r)
	})
}
