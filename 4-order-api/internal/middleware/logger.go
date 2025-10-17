package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger returns a new http.Handler that logs requests to the given logger.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Our custom ResponseWriter to capture status code
		lw := &loggingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     lw.statusCode,
			"duration":   duration.String(),
			"user_agent": r.UserAgent(),
			"remote_ip":  r.RemoteAddr,
		}).Info("Request handled")
	})
}

// loggingResponseWriter is a wrapper around http.ResponseWriter that captures the status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
