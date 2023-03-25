package db

import "time"

type CreateCollectionParams struct {
	Name         string    `json:"name" binding:"required"`
	Chain        int8      `json:"chain" binding:"required"`
	Address      string    `json:"address"`
	Creator      string    `json:"creator" binding:"required"`
	Type         int8      `json:"type" binding:"required"`
	Tax          int8      `json:"tax" binding:"required"`
	Currency     string    `json:"currency"`
	Visible      int8      `json:"visible"`
	Status       int8      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	Image        string    `json:"image" binding:"required"`
	Background   string    `json:"background" binding:"required"`
	Banner       string    `json:"banner" binding:"required"`
	Description  string    `json:"description"`
	Introduction string    `json:"introduction"`
	Properties   string    `json:"properties"`
	Twitter      string    `json:"twitter"`
	Instagram    string    `json:"instagram"`
	Discord      string    `json:"discord"`
	Web          string    `json:"web"`
}
