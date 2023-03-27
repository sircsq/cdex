package exchange

import (
	"bytes"
	"cdex/utils"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"sync"
	"time"
)

type OrderStatus string

const (
	PendingOrder  OrderStatus = "pending"
	FilledOrder   OrderStatus = "filled"
	CanceledOrder OrderStatus = "canceled"
)

type Match struct {
	Collection int     `json:"collection"`
	TokenID    int     `json:"token_id"`
	SizeFilled int     `json:"size_filled"`
	Price      float64 `json:"price"`
	Timestamp  int64   `json:"timestamp"`
	Ask        *Order2 `json:"ask"`
	Bid        *Order2 `json:"bid"`
}

type Order2 struct {
	ID    string
	Limit *Limit

	OrderRaw2
}

type Orders []*Order2

func (o Orders) Len() int {
	return len(o)
}

func (o Orders) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o Orders) Less(i, j int) bool {
	return o[i].Timestamp < o[j].Timestamp
}

type OrderRaw2 struct {
	Currency   string
	Owner      string
	Collection int
	TokenID    int
	Quantity   int
	Bid        bool
	Timestamp  int64
}

func (or *OrderRaw2) ID() string {
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

func NewOrder2(owner, currency string, bid bool, collection, tokenID, quantity int) *Order2 {
	raw := OrderRaw2{
		Currency:   currency,
		Owner:      owner,
		Collection: collection,
		TokenID:    tokenID,
		Quantity:   quantity,
		Bid:        bid,
		Timestamp:  time.Now().UnixNano(),
	}

	return &Order2{
		ID:        raw.ID(),
		OrderRaw2: raw,
	}
}

func (o *Order2) String() string {
	return fmt.Sprintf("[quantity:%v]", o.Quantity)
}

func (o *Order2) IsFilled() bool {
	return o.Quantity == 0
}

func (o *Order2) Type() string {
	if o.Bid {
		return "BID"
	}

	return "ASK"
}

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume int
}

type Limits []*Limit
type ByBestAsk struct {
	Limits
}

func (a ByBestAsk) Len() int {
	return len(a.Limits)
}

func (a ByBestAsk) Swap(i, j int) {
	a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i]
}

func (a ByBestAsk) Less(i, j int) bool {
	return a.Limits[i].Price < a.Limits[j].Price
}

type ByBestBid struct {
	Limits
}

func (b ByBestBid) Len() int {
	return len(b.Limits)
}

func (b ByBestBid) Swap(i, j int) {
	b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i]
}

func (b ByBestBid) Less(i, j int) bool {
	return b.Limits[i].Price > b.Limits[j].Price
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:       price,
		Orders:      []*Order2{},
		TotalVolume: 0,
	}
}

func (l *Limit) String() string {
	return fmt.Sprintf("[price: %v | Volume: %v | Orders: %v]", l.Price, l.TotalVolume, l.Orders)
}

func (l *Limit) AddOrder(o *Order2) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Quantity
}

func (l *Limit) DeleteOrder(o *Order2) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}

	o.Limit = nil
	l.TotalVolume -= o.Quantity

	sort.Sort(l.Orders)
}

func (l *Limit) Fill(o *Order2) []Match {
	var (
		matches        []Match
		ordersToDelete []*Order2
	)

	for _, order := range l.Orders {
		if o.IsFilled() {
			break
		}

		match, err := l.fillOrder(order, o)
		if err != nil {
			return []Match{}
		}

		matches = append(matches, match)
		l.TotalVolume -= match.SizeFilled
		if order.IsFilled() {
			ordersToDelete = append(ordersToDelete, order)
		}
	}

	for _, order := range ordersToDelete {
		l.DeleteOrder(order)
	}

	return matches
}

func (l *Limit) fillOrder(a, b *Order2) (Match, error) {
	if a.Collection != b.Collection || a.TokenID != b.TokenID {
		return Match{}, errors.New("collection or token not is not matched")
	}
	var (
		bid        *Order2
		ask        *Order2
		sizeFilled int
	)

	if a.Bid {
		bid = a
		ask = b
	} else {
		bid = b
		ask = a
	}

	if a.Quantity > b.Quantity {
		a.Quantity -= b.Quantity
		sizeFilled = b.Quantity
		b.Quantity = 0
	} else {
		b.Quantity -= a.Quantity
		sizeFilled = a.Quantity
		a.Quantity = 0
	}

	return Match{
		Collection: a.Collection,
		TokenID:    a.TokenID,
		SizeFilled: sizeFilled,
		Price:      l.Price,
		Timestamp:  time.Now().UnixNano(),
		Bid:        bid,
		Ask:        ask,
	}, nil
}

type OrderBook struct {
	mu        *sync.RWMutex
	asks      []*Limit
	bids      []*Limit
	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
	Orders    map[string]*Order2
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
		Orders:    make(map[string]*Order2),
		mu:        &sync.RWMutex{},
	}
}

func (ob *OrderBook) placeMarketOrder(o *Order2) ([]Match, error) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var (
		matches []Match
	)

	if o.Bid {
		if o.Quantity > ob.AskTotalVolume() {
			return matches, fmt.Errorf("not enough volume [quantity: %v] for market order [quantity: %v]", ob.AskTotalVolume(), o.Quantity)
		}

		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(false, limit)
			}
		}
	} else {
		if o.Quantity > ob.BidTotalVolume() {
			return matches, fmt.Errorf("not enough volume [quantity: %v] for market order [quantity: %v]", ob.BidTotalVolume(), o.Quantity)
		}
		for _, limit := range ob.Bids() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(true, limit)
			}
		}
	}

	return matches, nil
}

func (ob *OrderBook) placeLimitOrder(price float64, o *Order2) ([]Match, error) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var (
		limit            *Limit
		limitCounterpart *Limit
		matches          []Match
	)

	if o.Bid {
		limit = ob.BidLimits[price]
		limitCounterpart = ob.AskLimits[price]
	} else {
		limit = ob.AskLimits[price]
		limitCounterpart = ob.BidLimits[price]
	}

	if limitCounterpart == nil {
		if limit == nil {
			limit = NewLimit(price)
			if o.Bid {
				ob.bids = append(ob.bids, limit)
				ob.BidLimits[price] = limit
			} else {
				ob.asks = append(ob.asks, limit)
				ob.AskLimits[price] = limit
			}
		}

		logrus.WithFields(logrus.Fields{
			"price": limit.Price,
			"type":  o.Type(),
			"size":  o.Quantity,
			"owner": o.Owner,
		}).Info("new limit order")

		ob.Orders[o.ID] = o
		limit.AddOrder(o)
	} else {
		limitMatches := limitCounterpart.Fill(o)
		matches = append(matches, limitMatches...)

		if len(limitCounterpart.Orders) == 0 {
			if o.Bid {
				ob.clearLimit(false, limitCounterpart)
			} else {
				ob.clearLimit(true, limitCounterpart)
			}
		}
	}

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.bids = append(ob.bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.asks = append(ob.asks, limit)
			ob.AskLimits[price] = limit
		}
	}

	return matches, nil
}

func (ob *OrderBook) clearLimit(bid bool, l *Limit) {
	if bid {
		delete(ob.BidLimits, l.Price)
		for i := 0; i < len(ob.bids); i++ {
			if ob.bids[i] == l {
				ob.bids[i] = ob.bids[len(ob.bids)-1]
				ob.bids = ob.bids[:len(ob.bids)-1]
			}
		}
	} else {
		delete(ob.AskLimits, l.Price)
		for i := 0; i < len(ob.asks); i++ {
			if ob.asks[i] == l {
				ob.asks[i] = ob.asks[len(ob.asks)-1]
				ob.asks = ob.asks[:len(ob.asks)-1]
			}
		}
	}
}

func (ob *OrderBook) CancelOrder(o *Order2) {
	limit := o.Limit
	limit.DeleteOrder(o)
}

func (ob *OrderBook) BidTotalVolume() int {
	totalVolume := 0
	for i := 0; i < len(ob.bids); i++ {
		totalVolume += ob.bids[i].TotalVolume
	}
	return totalVolume
}

func (ob *OrderBook) AskTotalVolume() int {
	totalVolume := 0
	for i := 0; i < len(ob.asks); i++ {
		totalVolume += ob.asks[i].TotalVolume
	}
	return totalVolume
}

func (ob *OrderBook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

func (ob *OrderBook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}
