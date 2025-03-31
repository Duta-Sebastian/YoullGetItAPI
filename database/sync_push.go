package database

import (
	"YoullGetItAPI/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// PostCvSyncPushData updates or inserts CV data for a user
func PostCvSyncPushData(ctx context.Context, db *sql.DB, userId string, records []models.CvRecord) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	if len(records) == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	record := records[0]

	if record.LastChanged.IsZero() {
		record.LastChanged = time.Now()
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO CV (user_id, cv_data, last_changed)
		VALUES ($3, $1, $2)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
		cv_data = $1,
		last_changed = $2`,
		record.CvData, record.LastChanged, userId)

	if err != nil {
		return fmt.Errorf("failed to insert CV record: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// PostUserSyncPushData updates user data for a given user
func PostUserSyncPushData(ctx context.Context, db *sql.DB, userId string, records []models.UserRecord) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	if len(records) == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	record := records[0]

	if record.LastChanged.IsZero() {
		record.LastChanged = time.Now()
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE auth_user
		SET username=$1, last_changed=$2
		WHERE auth0_id=$3`,
		record.Username, record.LastChanged, userId)

	if err != nil {
		return fmt.Errorf("failed to update user record: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with auth0_id: %s", userId)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// PostJobCartSyncPushData updates or inserts job cart data for a user
func PostJobCartSyncPushData(ctx context.Context, db *sql.DB, userId string, records []models.JobRecord) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	if len(records) == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO job_cart (user_id, job_id, job_data, last_changed, status, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, job_id) 
		DO UPDATE SET 
			job_data = $3,
			last_changed = $4,
			status = $5,
			is_deleted = $6
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Error closing statement: %v", err)
		}
	}(stmt)

	for _, record := range records {
		if record.LastChanged.IsZero() {
			record.LastChanged = time.Now()
		}

		_, err = stmt.ExecContext(ctx,
			userId,
			record.JobId,
			record.JobData,
			record.LastChanged,
			record.Status,
			record.IsDeleted)

		if err != nil {
			return fmt.Errorf("failed to upsert job record (job_id: %s): %v", record.JobId, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
