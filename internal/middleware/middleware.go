package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging middleware logs all HTTP requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		log.Printf("%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
			r.RemoteAddr,
		)
	})
}

// CORS middleware adds CORS headers
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
		
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// APIKeyAuth middleware validates API keys
func APIKeyAuth(next http.Handler, validKeys []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, `{"error":"Missing X-API-Key header"}`, http.StatusUnauthorized)
			return
		}
		
		// Validate API key
		valid := false
		for _, key := range validKeys {
			if apiKey == key {
				valid = true
				break
			}
		}
		
		if !valid {
			http.Error(w, `{"error":"Invalid API key"}`, http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Recover middleware recovers from panics
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

