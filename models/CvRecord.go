package models

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type CvRecord struct {
	CvData      []byte    `json:"cv_data"`
	LastChanged time.Time `json:"last_changed"`
}

func (c *CvRecord) UnmarshalJSON(data []byte) error {
	type TempCvRecord struct {
		CvData      string    `json:"cv_data"`
		LastChanged time.Time `json:"last_changed"`
	}

	var temp TempCvRecord
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	decodedData, err := base64.StdEncoding.DecodeString(temp.CvData)
	if err != nil {
		return err
	}

	c.CvData = decodedData
	c.LastChanged = temp.LastChanged
	return nil
}
