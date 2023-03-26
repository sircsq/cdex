package api

import (
	"cdex/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type createCollectionRequest struct {
	Name         string `json:"name" binding:"required"`
	Chain        int8   `json:"chain" binding:"numeric"`
	Address      string `json:"address"`
	Creator      string `json:"creator" binding:"required"`
	Type         int8   `json:"type" binding:"numeric"`
	Tax          int8   `json:"tax" binding:"numeric"`
	Symbol       string `json:"symbol" binding:"required"`
	Currency     string `json:"currency" binding:"required"`
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
	var (
		err error
		req createCollectionRequest
		c   *db.Collection
	)
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
		Symbol:       req.Symbol,
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
	c, err = s.store.InsertCollection(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, c)
}

func (s *Server) listCollection(ctx *gin.Context) {
	var (
		err         error
		page        int64 = 1
		pageSize    int64 = 10
		collections []*db.Collection
	)

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

	collections, err = s.store.GetCollections(ctx, int(page), int(pageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, collections)
}

func (s *Server) listAddressCollection(ctx *gin.Context) {
	var (
		err         error
		address     string
		page        int64 = 1
		pageSize    int64 = 10
		collections []*db.Collection
	)
	address = ctx.Param("address")
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

	collections, err = s.store.GetCollectionByCreator(ctx, address, int(page), int(pageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, collections)
}
