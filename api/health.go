package api

import (
	"net/http"
)

// RegisterHealthRoute registers the public route
func RegisterHealthRoute(router *http.ServeMux) {
	router.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"The API is healthy"}`))
	}))
}
