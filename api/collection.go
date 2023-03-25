package api

import (
	"cdex/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type createCollectionRequest struct {
	Name         string `json:"name" binding:"required"`
	Chain        int8   `json:"chain" binding:"numeric"`
	Address      string `json:"address"`
	Creator      string `json:"creator" binding:"required"`
	Type         int8   `json:"type" binding:"numeric"`
	Tax          int8   `json:"tax" binding:"numeric"`
	Currency     string `json:"currency"`
	Visible      int8   `json:"visible"`
	Status       int8   `json:"status"`
	Image        string `json:"image" binding:"required"`
	Background   string `json:"background" binding:"required"`
	Banner       string `json:"banner" binding:"required"`
	Description  string `json:"description"`
	Introduction string `json:"introduction"`
	Properties   string `json:"properties"`
	Twitter      string `json:"twitter"`
	Instagram    string `json:"instagram"`
	Discord      string `json:"discord"`
	Web          string `json:"web"`
}

func (s *Server) createCollection(ctx *gin.Context) {
	var req createCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateCollectionParams{
		Name:         req.Name,
		Chain:        req.Chain,
		Address:      req.Address,
		Creator:      req.Creator,
		Type:         req.Type,
		Tax:          req.Tax,
		Currency:     req.Currency,
		Visible:      req.Visible,
		Status:       req.Status,
		CreatedAt:    time.Now(),
		Image:        req.Image,
		Background:   req.Background,
		Banner:       req.Banner,
		Description:  req.Description,
		Introduction: req.Introduction,
		Properties:   req.Properties,
		Twitter:      req.Twitter,
		Instagram:    req.Instagram,
		Discord:      req.Discord,
		Web:          req.Web,
	}
	err := s.store.Insert(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, db.Collection{})
}
