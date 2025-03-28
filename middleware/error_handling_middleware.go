package middleware

import (
	"log"
	"net/http"
)

func ErrorHandlingMiddleware(router *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and send a user-friendly message
				log.Printf("Error occurred: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		router.ServeHTTP(w, r)
	})
}
