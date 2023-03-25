package matching

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrder1 := NewOrder(true, 5)
	buyOrder2 := NewOrder(true, 1)
	buyOrder3 := NewOrder(true, 2)
	l.AddOrder(buyOrder1)
	l.AddOrder(buyOrder2)
	l.AddOrder(buyOrder3)
	l.DeleteOrder(buyOrder2)
	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderBook()
	sellOrder1 := NewOrder(false, 10)
	sellOrder2 := NewOrder(false, 5)
	ob.PlaceLimitOrder(10_000, sellOrder1)
	ob.PlaceLimitOrder(9_000, sellOrder2)
	assert(t, len(ob.asks), 2)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrder := NewOrder(false, 20)
	ob.PlaceLimitOrder(10_000, sellOrder)

	buyOrder := NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 10.0)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrder.IsFilled(), true)

	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderBook()
	buyOrder1 := NewOrder(true, 5)
	buyOrder2 := NewOrder(true, 8)
	buyOrder3 := NewOrder(true, 10)
	buyOrder4 := NewOrder(true, 1)

	ob.PlaceLimitOrder(5_000, buyOrder3)
	ob.PlaceLimitOrder(5_000, buyOrder4)
	ob.PlaceLimitOrder(9_000, buyOrder2)
	ob.PlaceLimitOrder(10_000, buyOrder1)

	assert(t, ob.BidTotalVolume(), 24.00)

	sellOrder := NewOrder(false, 20)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 4.0)
	assert(t, len(matches), 3)
	assert(t, len(ob.bids), 1)

	fmt.Printf("%+v\n", matches)
}

func TestPlaceMarketOrderMultiFill2(t *testing.T) {
	ob := NewOrderBook()
	sellOrder1 := NewOrder(false, 5)
	sellOrder2 := NewOrder(false, 8)
	sellOrder3 := NewOrder(false, 10)
	ob.PlaceLimitOrder(10_000, sellOrder1)
	ob.PlaceLimitOrder(9_000, sellOrder2)
	ob.PlaceLimitOrder(5_000, sellOrder3)

	assert(t, ob.AskTotalVolume(), 23.00)

	sellOrder := NewOrder(true, 20)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.AskTotalVolume(), 3.0)
	assert(t, len(matches), 3)
	assert(t, len(ob.asks), 1)

	fmt.Printf("%+v\n", matches)
}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderBook()
	buyOrder := NewOrder(true, 5)
	ob.PlaceLimitOrder(10_000.0, buyOrder)
	assert(t, ob.BidTotalVolume(), 5.0)
	assert(t, len(ob.bids), 1)
	ob.CancelOrder(buyOrder)
	assert(t, ob.BidTotalVolume(), 0.0)
}
