package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

// CreateUser adds a new user to the AUTH_USER table with proper error handling and context
func CreateUser(ctx context.Context, db *sql.DB, userID string) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	if userID == "" {
		return errors.New("userID cannot be empty")
	}

	query := `INSERT INTO AUTH_USER(auth0_id) VALUES ($1)`

	_, err := db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error inserting user (ID: %s): %v", userID, err)
		return err
	}

	return nil
}
