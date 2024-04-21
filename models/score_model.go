// models/score_model.go
package models

// Score represents the highscore data model
type Score struct {
	UserName string  `json:"user_name"`
	Country  string  `json:"country"`
	State    string  `json:"state"`
	Score    float64 `json:"score"`
}
