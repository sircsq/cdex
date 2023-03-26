package api

import (
	"cdex/exchange"
	"cdex/exchange/orderbook"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlaceOrderRequest struct {
	Market     exchange.Market    `json:"market" binding:"required"`
	Type       exchange.OrderType `json:"type" binding:"required"`
	Bid        bool               `json:"bid"`
	Collection int                `json:"collection" binding:"required,numeric"`
	TokenID    int                `json:"token_id" binding:"required,numeric"`
	Quantity   int                `json:"quantity" binding:"required,numeric"`
	Price      float64            `json:"price" binding:"required,numeric"`
}

type OrderData struct {
	Bid        bool    `json:"bid"`
	Collection int     `json:"collection"`
	TokenID    int     `json:"token_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Timestamp  int64   `json:"timestamp"`
}

type OrderBookData struct {
	Asks []*OrderData `json:"asks"`
	Bids []*OrderData `json:"bids"`
}

func (s *Server) getMartBook(ctx *gin.Context) {
	market := exchange.Market(ctx.Param("market"))
	ob, err := s.ex.OrderBook(market)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderBookData := OrderBookData{
		Asks: []*OrderData{},
		Bids: []*OrderData{},
	}

	for _, limit := range ob.Asks() {
		for _, o := range limit.Orders {
			order := OrderData{
				Price:      limit.Price,
				Bid:        o.Bid,
				Collection: o.Collection,
				TokenID:    o.TokenID,
				Quantity:   o.Quantity,
				Timestamp:  o.Timestamp,
			}
			orderBookData.Asks = append(orderBookData.Asks, &order)
		}
	}

	for _, limit := range ob.Bids() {
		for _, o := range limit.Orders {
			order := OrderData{
				Price:      limit.Price,
				Bid:        o.Bid,
				Collection: o.Collection,
				TokenID:    o.TokenID,
				Quantity:   o.Quantity,
				Timestamp:  o.Timestamp,
			}
			orderBookData.Bids = append(orderBookData.Bids, &order)
		}
	}

	ctx.JSON(http.StatusOK, orderBookData)
}

func (s *Server) placeOrder(ctx *gin.Context) {
	var (
		err error
		req PlaceOrderRequest
		ob  *orderbook.OrderBook
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, errorResponse(err))
		return
	}

	ob, err = s.ex.OrderBook(req.Market)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	o := orderbook.NewOrder(req.Bid, req.Collection, req.TokenID, req.Quantity)

	switch req.Type {
	case exchange.LimitOrder:
		ob.PlaceLimitOrder(req.Price, o)
		ctx.JSON(http.StatusOK, msgResponse("limit order placed"))
	case exchange.MarketOrder:
		_, err = ob.PlaceMarketOrder(o)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, msgResponse("market order placed"))
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("unknown order type")))
	}
}
