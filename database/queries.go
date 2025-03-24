package queries

import (
	"YoullGetItAPI/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"time"
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

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
		return nil, err
	}

	return db, nil
}

func CreateUser(db *sql.DB, userID string) error {
	query := `INSERT INTO AUTH_USER(auth0_id) VALUES ($1)`

	_, err := db.Exec(query, userID)
	if err != nil {
		log.Println("Error inserting user: ", err)
		return err
	}

	return nil
}

func GetSyncPullData(db *sql.DB, sinceTime *time.Time, userId string) ([]models.JobRecord, error) {
	rows, err := db.Query(`
				SELECT job_data, date_added, status
				FROM job_cart
				WHERE date_added > COALESCE($1, '1970-01-01'::timestamp) AND user_id = $2`, &sinceTime, userId)

	if err != nil {
		log.Println("Error getting sync data from database: ", err)
		return nil, err
	}

	var records []models.JobRecord
	for rows.Next() {
		var record models.JobRecord
		if err := rows.Scan(&record.JobData, &record.DateAdded, &record.Status); err != nil {
			log.Println("Error getting sync data from database: ", err)
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func PostSyncPushData(db *sql.DB, records []models.JobRecord, userId string) error {
	var valueStrings []string
	var args []interface{}

	for i, record := range records {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		args = append(args, userId, record.JobData, record.DateAdded, record.Status)
	}

	query := fmt.Sprintf(`
		INSERT INTO job_cart (user_id, job_data, date_added, status)
		VALUES %s
		ON CONFLICT (id) 
		DO UPDATE SET job_data = EXCLUDED.job_data, date_added = EXCLUDED.date_added, status = EXCLUDED.status`,
		strings.Join(valueStrings, ","),
	)

	// Execute the query with arguments
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Println("Error posting sync data from database: ", err)
		return fmt.Errorf("error performing bulk upsert: %v", err)
	}

	return nil
}
