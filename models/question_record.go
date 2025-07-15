package models

import "time"

type QuestionRecord struct {
	QuestionJSON         string    `json:"question_json"`
	LastChanged          time.Time `json:"last_changed"`
	IsShortQuestionnaire bool      `json:"is_short_questionnaire"`
}
