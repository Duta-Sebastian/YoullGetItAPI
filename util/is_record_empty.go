package util

import "YoullGetItAPI/models"

func IsRecordEmpty(records interface{}) bool {
	isEmpty := true
	if records != nil {
		switch v := records.(type) {
		case []models.JobRecord:
			isEmpty = len(v) == 0
		case []models.UserRecord:
			isEmpty = len(v) == 0
		case []models.CvRecord:
			isEmpty = len(v) == 0
		}
	}
	return isEmpty
}
