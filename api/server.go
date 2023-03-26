package api

import (
	"cdex/db"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Storage
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Storage) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.StaticFS("/static/", gin.Dir("./public/images", false))

	router.POST("/api/collection", server.createCollection)
	router.GET("/api/collection/list", server.listCollection)
	router.POST("/api/item", server.createItem)

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
