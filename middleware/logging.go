package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func LoggingMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Printf(
				"%s %s %s %v",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				time.Since(start),
			)
		})
	}
}
