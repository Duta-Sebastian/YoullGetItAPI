package queries

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

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
	if err != nil {
		log.Fatal("Error opening database: ", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
		return nil, err
	}

	return db, nil
}
