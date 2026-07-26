package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rwi *responseWriterInterceptor) WriteHeader(statusCode int) {
	rwi.statusCode = statusCode
	rwi.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rwi := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rwi, r)

		duration := time.Since(start)
		log.Printf("[%s] %s | Status: %d | Latency: %v", r.Method, r.URL.Path, rwi.statusCode, duration)
	})
}