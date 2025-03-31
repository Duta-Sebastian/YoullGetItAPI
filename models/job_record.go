package models

import "time"

type JobRecord struct {
	JobId       string    `json:"job_id"`
	JobData     string    `json:"job_data"`
	LastChanged time.Time `json:"last_changed"`
	Status      string    `json:"status"`
	IsDeleted   bool      `json:"is_deleted"`
}
