package db

import "time"

type Collection struct {
	ID           int       `json:id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Creator      string    `json:"creator"`
	Chain        int8      `json:"chain"`
	Visible      int8      `json:"visible"`
	Status       int8      `json:"status"`
	Type         int8      `json:"type"`
	Tax          int8      `json:"tax"`
	Symbol       string    `json:"symbol"`
	Currency     string    `json:"currency"`
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

type Item struct {
	Name        string    `json:"name"`
	Collection  int       `json:"collection"`
	TokenID     int       `json:"token_id"`
	Chain       int8      `json:"chain"`
	Creator     string    `json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Properties  string    `json:"properties"`
}
