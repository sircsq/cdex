package exchange

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

//func TestLimit(t *testing.T) {
//	l := NewLimit(10_000)
//	buyOrder1 := NewOrder(true, 1, 2, 5)
//	buyOrder2 := NewOrder(true, 1, 3, 1)
//	buyOrder3 := NewOrder(true, 1, 4, 2)
//	l.AddOrder(buyOrder1)
//	l.AddOrder(buyOrder2)
//	l.AddOrder(buyOrder3)
//	l.DeleteOrder(buyOrder2)
//	fmt.Println(l)
//}
//
//func TestPlaceLimitOrder(t *testing.T) {
//	ob := NewOrderBook()
//	sellOrder1 := NewOrder(false, 1, 2, 10)
//	sellOrder2 := NewOrder(false, 1, 3, 5)
//	ob.PlaceLimitOrder(10_000, sellOrder1)
//	ob.PlaceLimitOrder(9_000, sellOrder2)
//	assert(t, len(ob.asks), 2)
//}
//
//func TestPlaceMarketOrder(t *testing.T) {
//	ob := NewOrderBook()
//
//	sellOrder := NewOrder(false, 1, 2, 20)
//	ob.PlaceLimitOrder(10_000, sellOrder)
//
//	buyOrder := NewOrder(true, 1, 3, 10)
//	matches, _ := ob.PlaceMarketOrder(buyOrder)
//
//	assert(t, len(matches), 1)
//	assert(t, len(ob.asks), 1)
//	assert(t, ob.AskTotalVolume(), 10)
//	assert(t, matches[0].Ask, sellOrder)
//	assert(t, matches[0].Bid, buyOrder)
//	assert(t, matches[0].SizeFilled, 10)
//	assert(t, matches[0].Price, 10_000.0)
//	assert(t, buyOrder.IsFilled(), true)
//
//	fmt.Printf("%+v", matches)
//}
//
//func TestPlaceMarketOrderMultiFill(t *testing.T) {
//	ob := NewOrderBook()
//	buyOrder1 := NewOrder(true, 1, 2, 5)
//	buyOrder2 := NewOrder(true, 1, 3, 8)
//	buyOrder3 := NewOrder(true, 1, 4, 10)
//	buyOrder4 := NewOrder(true, 1, 5, 1)
//
//	ob.PlaceLimitOrder(5_000, buyOrder3)
//	ob.PlaceLimitOrder(5_000, buyOrder4)
//	ob.PlaceLimitOrder(9_000, buyOrder2)
//	ob.PlaceLimitOrder(10_000, buyOrder1)
//
//	assert(t, ob.BidTotalVolume(), 24)
//
//	sellOrder := NewOrder(false, 1, 2, 20)
//	matches, _ := ob.PlaceMarketOrder(sellOrder)
//
//	assert(t, ob.BidTotalVolume(), 4)
//	assert(t, len(matches), 3)
//	assert(t, len(ob.bids), 1)
//
//	fmt.Printf("%+v\n", matches)
//}
//
//func TestPlaceMarketOrderMultiFill2(t *testing.T) {
//	ob := NewOrderBook()
//	sellOrder1 := NewOrder(false, 1, 2, 5)
//	sellOrder2 := NewOrder(false, 1, 3, 8)
//	sellOrder3 := NewOrder(false, 1, 4, 10)
//	ob.PlaceLimitOrder(10_000, sellOrder1)
//	ob.PlaceLimitOrder(9_000, sellOrder2)
//	ob.PlaceLimitOrder(5_000, sellOrder3)
//
//	assert(t, ob.AskTotalVolume(), 23)
//
//	buyOrder := NewOrder(true, 1, 2, 20)
//	matches, _ := ob.PlaceMarketOrder(buyOrder)
//
//	assert(t, ob.AskTotalVolume(), 3)
//	assert(t, len(matches), 3)
//	assert(t, len(ob.asks), 1)
//
//	fmt.Printf("%+v\n", matches)
//}
//
//func TestCancelOrder(t *testing.T) {
//	ob := NewOrderBook()
//	buyOrder := NewOrder(true, 1, 2, 5)
//	ob.PlaceLimitOrder(10_000, buyOrder)
//	assert(t, ob.BidTotalVolume(), 5)
//	assert(t, len(ob.bids), 1)
//	ob.CancelOrder(buyOrder)
//	assert(t, ob.BidTotalVolume(), 0)
//}
