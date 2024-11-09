package api

import (
	db "blog-go-api/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllTagRequest struct {
	Name     string `form:"name"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) GetAllTag(ctx *gin.Context) {

	var req GetAllTagRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAllTagParams{
		Lower:  req.Name,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	tags, err := server.store.GetAllTag(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	count, err := server.store.CountAllTag(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := jsonResponseWithPaginate{
		jsonResponse: jsonResponse{
			Error:   false,
			Message: "successfully",
			Data:    tags,
		},
		Total: count,
	}

	ctx.JSON(http.StatusOK, payload)
}

type GetTagByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetTagById(ctx *gin.Context) {
	var req GetTagByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tag, err := server.store.GetTagById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    tag,
	})
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) CreateTag(ctx *gin.Context) {
	var req CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tag, err := server.store.CreateTag(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    tag,
	})
}

type UpdateTagRequest struct {
	ID   int64  `json:"id" binding:"required,min=1"`
	Name string `json:"name" binding:"required"`
}

func (server *Server) UpdateTag(ctx *gin.Context) {
	var req UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateTagParams{
		ID:   req.ID,
		Name: req.Name,
	}

	err := server.store.UpdateTag(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    nil,
	})
}

type DeleteTagRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) DeleteTag(ctx *gin.Context) {
	var req DeleteTagRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteTag(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    nil,
	})
}
