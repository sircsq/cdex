package db

import "time"

type Collection struct {
	ID           int64     `json:"ID"`
	Name         string    `json:"name"`
	Chain        int8      `json:"chain"`
	Address      string    `json:"address"`
	Creator      string    `json:"creator"`
	Type         int8      `json:"type"`
	Tax          int8      `json:"tax"`
	Currency     string    `json:"currency"`
	Visible      int8      `json:"visible"`
	Status       int8      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	Image        string    `json:"image"`
	Background   string    `json:"background"`
	Banner       string    `json:"banner"`
	Description  string    `json:"description"`
	Introduction string    `json:"introduction"`
	Properties   string    `json:"properties"`
	Twitter      string    `json:"twitter"`
	Instagram    string    `json:"instagram"`
	Discord      string    `json:"discord"`
	Web          string    `json:"web"`
}
