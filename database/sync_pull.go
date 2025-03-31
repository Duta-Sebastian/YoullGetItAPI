package database

import (
	"YoullGetItAPI/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func GetCvSyncPullData(ctx context.Context, db *sql.DB, userId string) ([]models.CvRecord, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	rows, err := db.QueryContext(ctx, `
		SELECT cv_data, last_changed
		FROM cv
		WHERE user_id = $1`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			fmt.Printf("close rows failed: %v", err)
		}
	}(rows)

	var cvRecords []models.CvRecord
	for rows.Next() {
		var cvRecord models.CvRecord
		if err = rows.Scan(&cvRecord.CvData, &cvRecord.LastChanged); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		cvRecords = append(cvRecords, cvRecord)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return cvRecords, nil
}

func GetUserSyncPullData(ctx context.Context, db *sql.DB, userId string) ([]models.UserRecord, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	rows, err := db.QueryContext(ctx, `
		SELECT username, last_changed
		FROM auth_user
		WHERE auth0_id = $1`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			fmt.Printf("close rows failed: %v", err)
		}
	}(rows)

	var userRecords []models.UserRecord
	for rows.Next() {
		var userRecord models.UserRecord
		var username sql.NullString
		var lastChanged sql.NullTime

		if err = rows.Scan(&username, &lastChanged); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		userRecord.Username = username.String

		if lastChanged.Valid {
			userRecord.LastChanged = lastChanged.Time
		} else {
			userRecord.LastChanged = time.Time{}
		}

		userRecords = append(userRecords, userRecord)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return userRecords, nil
}

func GetJobCartSyncPullData(ctx context.Context, db *sql.DB, userId string) ([]models.JobRecord, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	rows, err := db.QueryContext(ctx, `
		SELECT job_id, job_data, last_changed, status, is_deleted
		FROM job_cart
		WHERE user_id = $2`,
		userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			fmt.Printf("close rows failed: %v", err)
		}
	}(rows)

	var records []models.JobRecord
	for rows.Next() {
		var record models.JobRecord
		if err = rows.Scan(&record.JobId, &record.JobData, &record.LastChanged,
			&record.Status, &record.IsDeleted); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return records, nil
}
