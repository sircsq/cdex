package api

import (
	"cdex/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type createItemRequest struct {
	Name        string `json:"name"`
	Collection  int    `json:"collection" binding:"required"`
	TokenID     int    `json:"token_id" binding:"required"`
	Chain       int8   `json:"chain" binding:"numeric"`
	Creator     string `json:"creator" binding:"required"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Properties  string `json:"properties"`
}

func (s *Server) createItem(ctx *gin.Context) {
	var (
		err  error
		req  createItemRequest
		item *db.Item
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateItemParams{
		Name:        req.Name,
		Collection:  req.Collection,
		TokenID:     req.TokenID,
		Chain:       req.Chain,
		CreatedAt:   time.Now(),
		Creator:     req.Creator,
		Image:       req.Image,
		Description: req.Description,
		Properties:  req.Properties,
	}

	item, err = s.store.InsertItem(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (s *Server) listItem(ctx *gin.Context) {
	var (
		err      error
		page     int64 = 1
		pageSize int64 = 10
		items    []*db.Item
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

	items, err = s.store.GetItems(ctx, int(page), int(pageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (s *Server) listCollectionItem(ctx *gin.Context) {
	var (
		err      error
		id       int64
		page     int64 = 1
		pageSize int64 = 10
		items    []*db.Item
	)

	id, err = strconv.ParseInt(ctx.Param("collection"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
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

	items, err = s.store.GetCollectionItems(ctx, int(id), int(page), int(pageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}
