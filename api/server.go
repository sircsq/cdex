package api

import (
	"cdex/db"
	"cdex/exchange"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	ex      *exchange.Exchange
	store   db.Storage
	router  *gin.Engine
	mu      sync.RWMutex
	clients map[string]*websocket.Conn
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Storage) *Server {
	ex := exchange.NewExchange()
	server := &Server{
		ex:      ex,
		store:   store,
		clients: make(map[string]*websocket.Conn),
	}

	router := gin.Default()

	router.StaticFS("/static/", gin.Dir("./public/images", false))

	router.GET("/api/index/explore", server.listCollection)

	// collection
	router.POST("/api/collection", server.createCollection)
	router.GET("/api/collection/list", server.listCollection)
	router.GET("/api/collection/:address/list", server.listAddressCollection)

	// item
	router.POST("/api/item", server.createItem)
	router.GET("/api/item/list", server.listItem)
	router.GET("/api/item/:collection/list", server.listCollectionItem)

	// order
	router.POST("/api/order", server.createOrder)
	router.GET("/api/order/bids", server.getBidOrders)
	router.GET("/api/order/asks", server.getAskOrders)
	router.DELETE("/api/order/:id", server.cancelOrder)
	router.POST("/api/order/")

	server.router = router

	return server
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func msgResponse(msg string) gin.H {
	return gin.H{"msg": msg}
}
