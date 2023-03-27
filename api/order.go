package api

import (
	"cdex/exchange"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PlaceOrderRequest struct {
	Owner      string  `json:"owner" binding:"required"`
	Currency   string  `json:"currency" binding:"required"`
	Market     string  `json:"market" binding:"required"`
	Bid        int8    `json:"bid"`
	Collection int     `json:"collection" binding:"required,numeric"`
	TokenID    int     `json:"token_id" binding:"required,numeric"`
	Quantity   int     `json:"quantity" binding:"required,numeric"`
	Price      float64 `json:"price"`
}

type PlaceOrderRequest2 struct {
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
	ID         string  `json:"id"`
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

func (s *Server) getMarket(ctx *gin.Context) {

}

func (s *Server) getMartBook2(ctx *gin.Context) {
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

type PlaceOrderResponse2 struct {
	OrderID string           `json:"order_id"`
	Matches []exchange.Match `json:"matches"`
}

type PlaceOrderResponse struct {
	OrderID string `json:"order_id"`
}

func (s *Server) createOrder(ctx *gin.Context) {
	var (
		err error
		req PlaceOrderRequest
		res PlaceOrderResponse
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Quantity != 1 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid quantity")))
		return
	}

	order := exchange.NewOrder(req.Owner, req.Currency, req.Bid, req.Collection, req.TokenID, req.Quantity, req.Price)
	err = s.store.Insert(ctx, order)
	if err != nil {
		fmt.Println("--->>>", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res.OrderID = order.ID

	ctx.JSON(http.StatusOK, res)
}

func (s *Server) getBidOrders(ctx *gin.Context) {
	var (
		err      error
		page     int64 = 1
		pageSize int64 = 10
		orders   []*exchange.Order
		sort     = "desc"
		status   = "pending"
	)
	sortParam, _ := ctx.GetQuery("sort")
	if len(sortParam) != 0 {
		sort = sortParam
	}
	statusParam, _ := ctx.GetQuery("status")
	if len(statusParam) != 0 {
		status = statusParam
	}

	if pageStr, ok := ctx.GetQuery("page"); ok {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	if pageSizeStr, ok := ctx.GetQuery("pageSize"); ok {
		pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	orders, err = s.store.GetOrders(ctx, 1, int(page), int(pageSize), status, sort)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(len(orders))
	ctx.JSON(http.StatusOK, orders)
}

func (s *Server) getAskOrders(ctx *gin.Context) {
	var (
		err      error
		page     int64 = 1
		pageSize int64 = 10
		orders   []*exchange.Order
		sort     = "desc"
		status   = "pending"
	)
	sortParam, _ := ctx.GetQuery("sort")
	if len(sortParam) != 0 {
		sort = sortParam
	}
	statusParam, _ := ctx.GetQuery("status")
	if len(statusParam) != 0 {
		status = statusParam
	}

	if pageStr, ok := ctx.GetQuery("page"); ok {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	if pageSizeStr, ok := ctx.GetQuery("pageSize"); ok {
		pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	orders, err = s.store.GetOrders(ctx, 0, int(page), int(pageSize), status, sort)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (s *Server) placeOrder2(ctx *gin.Context) {
	var (
		err error
		req PlaceOrderRequest2
		res PlaceOrderResponse2
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Quantity != 1 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid quantity")))
		return
	}

	order := exchange.NewOrder2(req.Owner, req.Currency, req.Bid, req.Collection, req.TokenID, req.Quantity)
	res.OrderID = order.ID

	switch req.Type {
	case exchange.LimitOrder:
		res.Matches, err = s.ex.PlaceLimitOrder(req.Market, req.Price, order)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, res)
	case exchange.MarketOrder:
		res.Matches, err = s.ex.PlaceMarketOrder(req.Market, order)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		ctx.JSON(http.StatusOK, res)
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("unknown order type")))
	}
}

type CancelOrderRequest struct {
	Market exchange.Market `json:"market" binding:"required"`
	Bid    bool            `json:"bid" binding:"required"`
	ID     string          `json:"id" binding:"required"`
}

func (s *Server) cancelOrder(ctx *gin.Context) {
	var err error
	orderID := ctx.Param("id")

	err = s.store.UpdateOrderStatus(ctx, orderID, "canceled")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, msgResponse("order canceled"))
}

func (s *Server) cancelOrder2(ctx *gin.Context) {
	var (
		err error
		req CancelOrderRequest
		ob  *exchange.OrderBook
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
