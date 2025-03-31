package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

// DBConnections holds all database connections used by the application
type DBConnections struct {
	UserDB       *sql.DB
	AuthActionDB *sql.DB
}

// SetupDatabases initializes all database connections with proper configuration
func SetupDatabases() (*DBConnections, error) {
	userDb, err := GetDBConnection("app_user")
	if err != nil {
		return nil, err
	}
	configureConnectionPool(userDb)

	authActionDb, err := GetDBConnection("auth_action")
	if err != nil {
		return nil, err
	}
	configureConnectionPool(authActionDb)

	return &DBConnections{
		UserDB:       userDb,
		AuthActionDB: authActionDb,
	}, nil
}

// configureConnectionPool sets optimal connection pool parameters
func configureConnectionPool(db *sql.DB) {
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("Warning: Error pinging database with user: %v", err)
	}
}

// CloseConnections properly closes all database connections
func (dbc *DBConnections) CloseConnections() {
	if dbc.UserDB != nil {
		if err := dbc.UserDB.Close(); err != nil {
			log.Printf("Error closing user database connection: %v", err)
		}
	}

	if dbc.AuthActionDB != nil {
		if err := dbc.AuthActionDB.Close(); err != nil {
			log.Printf("Error closing auth action database connection: %v", err)
		}
	}
}

// GetDBConnection initializes the database connection with different roles based on who the user is
func GetDBConnection(user string) (*sql.DB, error) {
	var dbUser, dbPassword, dbName, dbHost, dbPort string

	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")

	switch user {
	case "app_user":
		dbUser = os.Getenv("APP_USER")
		dbPassword = os.Getenv("APP_USER_PASSWORD")
	case "auth_action":
		dbUser = os.Getenv("AUTH0_ACTION_USER")
		dbPassword = os.Getenv("AUTH0_ACTION_USER_PASSWORD")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=require",
		dbUser, dbPassword, dbName, dbHost, dbPort)
	db, err := sql.Open("postgres", connStr)
	return db, err
}
