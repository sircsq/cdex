package api

import (
	"cdex/exchange"
	"cdex/exchange/orderbook"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlaceOrderRequest struct {
	Owner      string             `json:"owner" binding:"required"`
	Currency   string             `json:"currency" binding:"required"`
	Market     exchange.Market    `json:"market" binding:"required"`
	Type       exchange.OrderType `json:"type" binding:"required"`
	Bid        bool               `json:"bid"`
	Collection int                `json:"collection" binding:"required,numeric"`
	TokenID    int                `json:"token_id" binding:"required,numeric"`
	Quantity   int                `json:"quantity" binding:"required,numeric"`
	Price      float64            `json:"price"`
}

type OrderData struct {
	ID         int64   `json:"id"`
	Currency   string  `json:"currency"`
	Owner      string  `json:"owner"`
	Bid        bool    `json:"bid"`
	Collection int     `json:"collection"`
	TokenID    int     `json:"token_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Timestamp  int64   `json:"timestamp"`
}

type OrderBookData struct {
	TotalBidVolume int          `json:"total_bid_volume"`
	TotalAskVolume int          `json:"total_ask_volume"`
	Asks           []*OrderData `json:"asks"`
	Bids           []*OrderData `json:"bids"`
}

func (s *Server) getMartBook(ctx *gin.Context) {
	market := exchange.Market(ctx.Param("market"))
	ob, err := s.ex.OrderBook(market)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderBookData := OrderBookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*OrderData{},
		Bids:           []*OrderData{},
	}

	for _, limit := range ob.Asks() {
		for _, o := range limit.Orders {
			order := OrderData{
				ID:         o.ID,
				Currency:   o.Currency,
				Owner:      o.Owner,
				Bid:        o.Bid,
				Collection: o.Collection,
				TokenID:    o.TokenID,
				Quantity:   o.Quantity,
				Timestamp:  o.Timestamp,
				Price:      limit.Price,
			}
			orderBookData.Asks = append(orderBookData.Asks, &order)
		}
	}

	for _, limit := range ob.Bids() {
		for _, o := range limit.Orders {
			order := OrderData{
				ID:         o.ID,
				Currency:   o.Currency,
				Owner:      o.Owner,
				Bid:        o.Bid,
				Collection: o.Collection,
				TokenID:    o.TokenID,
				Quantity:   o.Quantity,
				Timestamp:  o.Timestamp,
				Price:      limit.Price,
			}
			orderBookData.Bids = append(orderBookData.Bids, &order)
		}
	}

	ctx.JSON(http.StatusOK, orderBookData)
}

func (s *Server) placeOrder(ctx *gin.Context) {
	var (
		err     error
		req     PlaceOrderRequest
		ob      *orderbook.OrderBook
		matches []orderbook.Match
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

	o := orderbook.NewOrder(req.Owner, req.Currency, req.Bid, req.Collection, req.TokenID, req.Quantity)

	switch req.Type {
	case exchange.LimitOrder:
		ob.PlaceLimitOrder(req.Price, o)
		ctx.JSON(http.StatusOK, msgResponse("limit order placed"))
	case exchange.MarketOrder:
		matches, err = ob.PlaceMarketOrder(o)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, matches)
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("unknown order type")))
	}
}

type CancelOrderRequest struct {
	Market exchange.Market `json:"market" binding:"required"`
	Bid    bool            `json:"bid" binding:"required"`
	ID     int64           `json:"id" binding:"required"`
}

func (s *Server) cancelOrder(ctx *gin.Context) {
	var (
		err error
		req CancelOrderRequest
		ob  *orderbook.OrderBook
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ob, err = s.ex.OrderBook(req.Market)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Bid {
		for _, limit := range ob.Bids() {
			for _, o := range limit.Orders {
				if o.ID == req.ID {
					ob.CancelOrder(o)
					ctx.JSON(http.StatusOK, msgResponse("order canceled"))
					return
				}
			}
		}
	} else {
		for _, limit := range ob.Asks() {
			for _, o := range limit.Orders {
				if o.ID == req.ID {
					ob.CancelOrder(o)
					ctx.JSON(http.StatusOK, msgResponse("order canceled"))
					return
				}
			}
		}
	}

	ctx.JSON(http.StatusBadRequest, msgResponse("oder not found"))
}
