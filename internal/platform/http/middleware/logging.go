package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"appgents/internal/platform/logging"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// RequestLogger creates a middleware that logs request and response details
func RequestLogger(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Read and restore the request body
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Create response wrapper to capture response
			rw := &responseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
			}

			// Process request
			next.ServeHTTP(rw, r)

			// Log request/response details
			logger.Info("http_request", map[string]interface{}{
				"operation":      r.Method + " " + r.URL.Path,
				"method":         r.Method,
				"path":           r.URL.Path,
				"status":         rw.status,
				"duration_ms":    time.Since(start).Milliseconds(),
				"request_body":   string(bodyBytes),
				"response_body":  rw.body.String(),
				"content_length": r.ContentLength,
				"remote_addr":    r.RemoteAddr,
				"user_agent":     r.UserAgent(),
			},
			)
		})
	}
}
