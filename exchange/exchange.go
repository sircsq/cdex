package exchange

import (
	"cdex/exchange/orderbook"
	"errors"
)

type OrderType string

const (
	MarketOrder OrderType = "market"
	LimitOrder  OrderType = "limit"
)

type Market string

const (
	MarketFRA Market = "fra"
)

type Exchange struct {
	orderBooks map[Market]*orderbook.OrderBook
}

func NewExchange() *Exchange {
	orderBooks := make(map[Market]*orderbook.OrderBook)
	orderBooks[MarketFRA] = orderbook.NewOrderBook()

	return &Exchange{orderBooks: orderBooks}
}

func (ex *Exchange) OrderBook(market Market) (*orderbook.OrderBook, error) {
	ob, ok := ex.orderBooks[market]
	if !ok {
		return nil, errors.New("market not found")
	}

	return ob, nil
}
