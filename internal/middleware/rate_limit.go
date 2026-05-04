package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type ipEntry struct {
	count    int
	windowAt time.Time
	mu       sync.Mutex
}

var ipMap sync.Map

// RateLimit allows up to limit requests per minute per IP.
func RateLimit(limit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			val, _ := ipMap.LoadOrStore(ip, &ipEntry{windowAt: time.Now()})
			entry := val.(*ipEntry)

			entry.mu.Lock()
			if time.Since(entry.windowAt) > time.Minute {
				entry.count = 0
				entry.windowAt = time.Now()
			}
			entry.count++
			over := entry.count > limit
			entry.mu.Unlock()

			if over {
				w.Header().Set("Retry-After", "60")
				http.Error(w, `{"error":"too many requests"}`, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
