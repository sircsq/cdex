package exchange

import (
	"errors"
	"sync"
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
	Orders     map[string][]*Order2 // user => []*Order
	orderBooks map[Market]*OrderBook
	mu         *sync.RWMutex
}

func NewExchange() *Exchange {
	orderBooks := make(map[Market]*OrderBook)
	orderBooks[MarketFRA] = NewOrderBook()

	return &Exchange{
		orderBooks: orderBooks,
		Orders:     make(map[string][]*Order2),
		mu:         &sync.RWMutex{},
	}
}

func (ex *Exchange) OrderBook(market Market) (*OrderBook, error) {
	ob, ok := ex.orderBooks[market]
	if !ok {
		return nil, errors.New("market not found")
	}

	return ob, nil
}

func (ex *Exchange) PlaceLimitOrder(market Market, price float64, order *Order2) ([]Match, error) {
	var (
		err     error
		matches []Match
	)
	ob, err := ex.OrderBook(market)
	if err != nil {
		return matches, err
	}

	matches, err = ob.placeLimitOrder(price, order)
	if err != nil {
		return matches, err
	}

	ex.mu.Lock()
	if len(matches) == 0 {
		ex.Orders[order.Owner] = append(ex.Orders[order.Owner], order)
	} else {
		for _, m := range matches {
			delete(ex.Orders, m.Ask.Owner)
			delete(ex.Orders, m.Bid.Owner)
		}
	}
	ex.mu.Unlock()

	return matches, nil
}

func (ex *Exchange) PlaceMarketOrder(market Market, order *Order2) ([]Match, error) {
	var (
		err     error
		matches []Match
	)

	ob := ex.orderBooks[market]

	matches, err = ob.placeMarketOrder(order)
	if err != nil {
		return matches, err
	}

	ex.mu.Lock()
	if len(matches) > 0 {
		for _, m := range matches {
			delete(ex.Orders, m.Ask.Owner)
			delete(ex.Orders, m.Bid.Owner)
		}
	}
	ex.mu.Unlock()

	return matches, nil
}
