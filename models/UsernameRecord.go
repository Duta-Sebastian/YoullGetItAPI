package models

import "time"

type UsernameRecord struct {
	Username    string    `json:"username"`
	LastChanged time.Time `json:"last_changed"`
}
