package api

import (
	"cdex/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createItemRequest struct {
	Name         string `json:"name"`
	CollectionID int    `json:"collection_id" binding:"required"`
	TokenID      int    `json:"token_id" binding:"required"`
	Chain        int8   `json:"chain" binding:"required"`
	Creator      string `json:"creator" binding:"required"`
	Image        string `json:"image"`
	Description  string `json:"description"`
	Properties   string `json:"properties"`
}

func (s *Server) createItem(ctx *gin.Context) {
	//var req createItemRequest
	//if err := ctx.ShouldBindJSON(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	//}
	//arg := db.CreateItemParams{
	//	Name:         req.Name,
	//	CollectionID: req.CollectionID,
	//	TokenID:      req.TokenID,
	//	Chain:        req.Chain,
	//	CreatedAt:    time.Now(),
	//	Creator:      req.Creator,
	//	Image:        req.Image,
	//	Description:  req.Description,
	//	Properties:   req.Properties,
	//}
	//
	//err := s.store.Insert(ctx, arg)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//}

	ctx.JSON(http.StatusOK, db.Item{})
}
