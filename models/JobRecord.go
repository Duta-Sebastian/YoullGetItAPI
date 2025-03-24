package models

import "time"

type JobRecord struct {
	JobData   string    `json:"job_data"`
	DateAdded time.Time `json:"date_added"`
	Status    string    `json:"status"`
}
