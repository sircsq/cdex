package db

import "time"

type CreateCollectionParams struct {
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

// name varchar(32) not null,
// collection_id integer not null,
// token_id integer,
// chain smallint not null,
// creator varchar(65) not null,
// created_at timestamp not null,
// image varchar(128) not null,
// description varchar(128),
// properties varchar(128),
type CreateItemParams struct {
	Name         string    `json:"name"`
	CollectionID int       `json:"collection_id"`
	TokenID      int       `json:"token_id"`
	Chain        int8      `json:"chain"`
	CreatedAt    time.Time `json:"created_at"`
	Creator      string    `json:"creator"`
	Image        string    `json:"image"`
	Description  string    `json:"description"`
	Properties   string    `json:"properties"`
}
