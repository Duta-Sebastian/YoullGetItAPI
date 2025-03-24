package api

import (
	"net/http"
)

// RegisterPublicRoute registers the public route
func RegisterPublicRoute(router *http.ServeMux) {
	router.Handle("/api/public", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello from a public endpoint! You don't need to be authenticated to see this."}`))
	}))
}
