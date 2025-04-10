package main

import (
	"YoullGetItAPI/api"
	"YoullGetItAPI/database"
	"YoullGetItAPI/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading the .env file: %v\n", err)
	}

	dbConnections, err := database.SetupDatabases()
	if err != nil {
		log.Fatalf("Failed to set up database connections: %v", err)
	}
	defer dbConnections.CloseConnections()

	router := http.NewServeMux()

	api.RegisterHealthRoute(router)
	api.RegisterCreateUserRoute(router, dbConnections.AuthActionDB)
	api.RegisterSyncPushRoutes(router, dbConnections.UserDB)
	api.RegisterSyncPullRoute(router, dbConnections.UserDB)

	handler := middleware.ChainMiddleware(router, middleware.LoggingMiddleware, middleware.ErrorHandlingMiddleware)

	log.Print("Server listening on http://localhost:3010")
	if err := http.ListenAndServe("0.0.0.0:3010", handler); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
