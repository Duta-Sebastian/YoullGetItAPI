package queries

import (
	"database/sql"
	"log"
)

func CreateUser(db *sql.DB, userID string) error {
	query := `INSERT INTO AUTH_USER(auth0_id) VALUES ($1)`

	_, err := db.Exec(query, userID)
	if err != nil {
		log.Println("Error inserting user: ", err)
		return err
	}

	return nil
}
