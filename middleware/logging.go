package middleware

import (
	"net/http"
	"time"
	"go.uber.org/zap"
)


func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Log the incoming request
		log := GetLogger()
		log.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		duration := time.Since(start) // Suggestions from AI
		
		// Log the response
		log.Info("HTTP Response",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status_code", wrapped.statusCode),
			zap.Duration("duration", duration),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
} 