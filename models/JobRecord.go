package models

import "time"

type JobRecord struct {
	JobId     string    `json:"job_id"`
	JobData   string    `json:"job_data"`
	DateAdded time.Time `json:"date_added"`
	Status    string    `json:"status"`
}
