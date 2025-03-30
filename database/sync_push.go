package queries

import (
	"YoullGetItAPI/models"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

func PostCvSyncPushData(db *sql.DB, userId string, records []models.CvRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	record := records[0]

	if record.LastChanged.IsZero() {
		record.LastChanged = time.Now()
	}

	_, err = tx.Exec(`
		UPDATE CV
		SET cv_data = $1, last_changed = $2
		WHERE user_id = $3`,
		record.CvData, record.LastChanged, userId)

	if err != nil {
		return fmt.Errorf("failed to insert CV record: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func PostUserSyncPushData(db *sql.DB, userId string, records []models.UserRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	record := records[0]

	if record.LastChanged.IsZero() {
		record.LastChanged = time.Now()
	}

	_, err = tx.Exec(`
		UPDATE auth_user
		SET username=$1, last_changed=$2
		WHERE user_id=$3`,
		record.Username, record.LastChanged, userId)

	if err != nil {
		return fmt.Errorf("failed to update user record: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
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
		ON CONFLICT (user_id, job_id)
		DO UPDATE SET job_data = EXCLUDED.job_data, date_added = EXCLUDED.date_added, status = EXCLUDED.status`,
		strings.Join(valueStrings, ","),
	)

	_, err := db.Exec(query, args...)
	if err != nil {
		log.Println("Error posting sync data from database: ", err)
		return fmt.Errorf("error performing bulk upsert: %v", err)
	}

	return nil
}
