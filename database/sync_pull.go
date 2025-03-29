package queries

import (
	"YoullGetItAPI/models"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func GetCvSyncPullData(db *sql.DB, userId string) ([]models.CvRecord, error) {
	rows, err := db.Query(`
				SELECT cv_data, last_changed
				FROM cv
				WHERE user_id = $1`, userId)

	if err != nil {
		return nil, err
	}

	var CvRecords []models.CvRecord
	for rows.Next() {
		var cvRecord models.CvRecord
		if err := rows.Scan(&cvRecord.CvData, &cvRecord.LastChanged); err != nil {
			return nil, err
		}
		CvRecords = append(CvRecords, cvRecord)
	}

	return CvRecords, nil
}

func GetUserSyncPullData(db *sql.DB, userId string) ([]models.UserRecord, error) {
	rows, err := db.Query(`
				SELECT username, last_changed
				FROM auth_user
				WHERE auth0_id = $1`, userId)

	if err != nil {
		return nil, err
	}

	var userRecords []models.UserRecord
	for rows.Next() {
		var userRecord models.UserRecord
		if err := rows.Scan(&userRecord.Username, &userRecord.LastChanged); err != nil {
			return nil, err
		}
		userRecords = append(userRecords, userRecord)
	}

	return userRecords, nil
}

func GetSyncPullData(db *sql.DB, sinceTime *time.Time, userId string) ([]models.JobRecord, error) {
	rows, err := db.Query(`
				SELECT job_id, job_data, date_added, status
				FROM job_cart
				WHERE date_added > COALESCE($1, '1970-01-01'::timestamp) AND user_id = $2`, &sinceTime, userId)

	if err != nil {
		log.Println("Error getting sync data from database: ", err)
		return nil, err
	}

	var records []models.JobRecord
	for rows.Next() {
		var record models.JobRecord
		if err := rows.Scan(&record.JobId, &record.JobData, &record.DateAdded, &record.Status); err != nil {
			log.Println("Error getting sync data from database: ", err)
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
