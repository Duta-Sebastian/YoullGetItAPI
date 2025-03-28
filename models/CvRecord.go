package models

import "time"

type CvRecord struct {
	CvData      []byte    `json:"cv_data"`
	LastChanged time.Time `json:"last_changed"`
}
