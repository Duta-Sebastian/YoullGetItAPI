package main

import (
	"YoullGetItAPI/api"
	"YoullGetItAPI/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	router := http.NewServeMux()

	api.RegisterPublicRoute(router)
	api.RegisterCreateUserRoute(router)
	api.RegisterSyncPushRoutes(router)
	api.RegisterSyncPullRoute(router)

	handler := middleware.ChainMiddleware(router, middleware.LoggingMiddleware, middleware.ErrorHandlingMiddleware)

	log.Print("Server listening on http://localhost:3010")
	if err := http.ListenAndServe("0.0.0.0:3010", handler); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
