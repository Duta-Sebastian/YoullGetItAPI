package models

import "time"

type UserRecord struct {
	Username    string    `json:"username"`
	LastChanged time.Time `json:"last_changed"`
}
