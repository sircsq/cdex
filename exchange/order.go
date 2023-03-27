package exchange

import (
	"bytes"
	"cdex/utils"
	"encoding/gob"
	"time"
)

type Order struct {
	ID         string      `json:"id"`
	Status     OrderStatus `json:"status"`
	Currency   string      `json:"currency"`
	Owner      string      `json:"owner"`
	Collection int         `json:"collection"`
	TokenID    int         `json:"token_id"`
	Quantity   int         `json:"quantity"`
	Bid        int8        `json:"bid"`
	Price      float64     `json:"price"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderRaw struct {
	Currency   string    `json:"currency"`
	Owner      string    `json:"owner"`
	Collection int       `json:"collection"`
	TokenID    int       `json:"token_id"`
	Quantity   int       `json:"quantity"`
	Bid        int8      `json:"bid"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
}

func (or *OrderRaw) ID() string {
	var (
		err error
		buf bytes.Buffer
	)

	enc := gob.NewEncoder(&buf)
	err = enc.Encode(*or)
	if err != nil {
		return ""
	}

	return utils.MD5(buf.Bytes())
}

func NewOrder(owner, currency string, bid int8, collection, tokenID, quantity int, price float64) *Order {
	raw := OrderRaw{
		Currency:   currency,
		Owner:      owner,
		Collection: collection,
		TokenID:    tokenID,
		Quantity:   quantity,
		Bid:        bid,
		Price:      price,
		CreatedAt:  time.Now(),
	}

	return &Order{
		ID:         raw.ID(),
		Status:     PendingOrder,
		Currency:   currency,
		Owner:      owner,
		Collection: collection,
		TokenID:    tokenID,
		Quantity:   quantity,
		Bid:        bid,
		Price:      price,
		CreatedAt:  time.Now(),
	}
}
