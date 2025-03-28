package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(router *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)

		router.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}
