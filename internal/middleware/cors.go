package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS applies restrictive cross-origin headers. Allowed origins are read from
// the CORS_ORIGINS env var (comma-separated). Defaults to localhost:3000 if unset.
func CORS(next http.Handler) http.Handler {
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		raw = "http://localhost:3000"
	}

	allowed := make(map[string]struct{})
	for _, o := range strings.Split(raw, ",") {
		allowed[strings.TrimSpace(o)] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowed[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "86400")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
