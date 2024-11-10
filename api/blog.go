package api

import (
	"fmt"
	"net/http"
	"os"

	db "blog-go-api/db/sqlc"
	"blog-go-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetAllBlogRequest struct {
	Name     string `form:"name"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=10"`
}

type GetAllBlogResponse struct {
	db.GetAllBlogRow
	BlogTags []db.GetBlogTagByBlogIdRow `json:"blog_tags"`
}

// GetAllBlog godoc
//
//	@Summary		Get All Blog
//	@Description	Get All Blog
//	@Tags			Blog
//	@Accept			json
//	@Produce		json
//	@Param			name		query		string	false	"Blog Name"
//	@Param			page_id		query		int		true	"Page ID"
//	@Param			page_size	query		int		true	"Page Size"
//	@Success		200			{object}	jsonResponseWithPaginate
//	@Router			/api/blog [get]
func (server *Server) GetAllBlog(ctx *gin.Context) {

	var req GetAllBlogRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAllBlogParams{
		Lower:  req.Name,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	var getAllBlogResponse []GetAllBlogResponse

	blogs, err := server.store.GetAllBlog(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, blog := range blogs {
		blogTags, err := server.store.GetBlogTagByBlogId(ctx, blog.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		getAllBlogResponse = append(getAllBlogResponse, GetAllBlogResponse{
			GetAllBlogRow: blog,
			BlogTags:      blogTags,
		})
	}

	count, err := server.store.CountAllBlog(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := jsonResponseWithPaginate{
		jsonResponse: jsonResponse{
			Error:   false,
			Message: "successfully",
			Data:    getAllBlogResponse,
		},
		Total: count,
	}

	ctx.JSON(http.StatusOK, payload)
}

type GetAllBlogWithTagRequest struct {
	Title    string `form:"title"`
	Tag      string `form:"tag"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=10"`
}

type GetAllBlogWithTagResponse struct {
	db.GetAllBlogWithTagRow
	BlogTags []db.GetBlogTagByBlogIdRow `json:"blog_tags"`
}

// GetAllBlogWithTag godoc
//
//	@Summary		Get All Blog With Tag
//	@Description	Get All Blog With Tag
//	@Tags			Blog
//	@Accept			json
//	@Produce		json
//	@Param			title		query		string	false	"Blog Title"
//	@Param			tag			query		string	false	"Tag Name"
//	@Param			page_id		query		int		true	"Page ID"
//	@Param			page_size	query		int		true	"Page Size"
//	@Success		200			{object}	jsonResponseWithPaginate
//	@Router			/api/blog/tag [get]
func (server *Server) GetAllBlogWithTag(ctx *gin.Context) {

	var req GetAllBlogWithTagRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAllBlogWithTagParams{
		Lower:   req.Title,
		Lower_2: req.Tag,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	var getAllBlogResponse []GetAllBlogWithTagResponse

	blogs, err := server.store.GetAllBlogWithTag(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, blog := range blogs {
		blogTags, err := server.store.GetBlogTagByBlogId(ctx, blog.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		getAllBlogResponse = append(getAllBlogResponse, GetAllBlogWithTagResponse{
			GetAllBlogWithTagRow: blog,
			BlogTags:             blogTags,
		})
	}

	argCount := db.CountAllBlogWithTagParams{
		Lower:   req.Title,
		Lower_2: req.Tag,
	}

	count, err := server.store.CountAllBlogWithTag(ctx, argCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := jsonResponseWithPaginate{
		jsonResponse: jsonResponse{
			Error:   false,
			Message: "successfully",
			Data:    getAllBlogResponse,
		},
		Total: count,
	}

	ctx.JSON(http.StatusOK, payload)
}

type GetBlogByUrlRequest struct {
	URL string `uri:"url" binding:"required,min=1"`
}

type GetBlogByUrlResponse struct {
	db.GetBlogByUrlRow
	BlogTags []db.GetBlogTagByBlogIdRow `json:"blog_tags"`
}

// GetBlogByUrl godoc
//
//	@Summary		Get Blog By URL
//	@Description	Get Blog By URL
//	@Tags			Blog
//	@Accept			json
//	@Produce		json
//	@Param			url	path		string	true	"Blog URL"
//	@Success		200	{object}	jsonResponse
//	@Router			/api/blog/{url} [get]
func (server *Server) GetBlogByUrl(ctx *gin.Context) {
	var req GetBlogByUrlRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get Blog by URL
	blog, err := server.store.GetBlogByUrl(ctx, req.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get Blog Tags
	blogTags, err := server.store.GetBlogTagByBlogId(ctx, blog.ID)

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data: GetBlogByUrlResponse{
			GetBlogByUrlRow: blog,
			BlogTags:        blogTags,
		},
	})
}

type GetBlogByIDRequest struct {
	ID int64 `form:"id" binding:"required,min=1"`
}

type GetBlogByIDResponse struct {
	db.GetBlogByIdRow
	BlogTags []db.GetBlogTagByBlogIdRow `json:"blog_tags"`
}

// GetBlogByID godoc
//
//	@Summary		Get Blog By ID
//	@Description	Get Blog By ID
//	@Tags			Blog
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int	true	"Blog ID"
//	@Success		200	{object}	jsonResponse
//	@Router			/api/blog/id [get]
//	@Security		BearerAuth
func (server *Server) GetBlogByID(ctx *gin.Context) {
	var req GetBlogByIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get Blog by ID
	blog, err := server.store.GetBlogById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get Blog Tags
	blogTags, err := server.store.GetBlogTagByBlogId(ctx, blog.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data: GetBlogByIDResponse{
			GetBlogByIdRow: blog,
			BlogTags:       blogTags,
		},
	})
}

type BlogTagRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	BlogID  int64 `json:"blog_id" binding:"required,min=1"`
	TagId   int64 `json:"tag_id" binding:"required,min=1"`
	Deleted bool  `json:"deleted"`
}

type CreateBlogRequest struct {
	Title    string           `json:"title" binding:"required"`
	Content  string           `json:"content" binding:"required"`
	Image    string           `json:"image" binding:"required"`
	URL      string           `json:"url" binding:"required"`
	BlogTags []BlogTagRequest `json:"blog_tags"`
}

type CreateBlogByIdResponse struct {
	db.Blog
	BlogTags []db.GetBlogTagByBlogIdRow `json:"blog_tags"`
}

// CreateBlog godoc
//
//	@Summary		Create Blog
//	@Description	Create Blog
//	@Tags			Blog
//	@Accept			json
//	@Produce		json
//	@Param			input	body		CreateBlogRequest	true	"Create information"
//	@Success		200		{object}	jsonResponse
//	@Router			/api/blog [post]
//	@Security		BearerAuth
func (server *Server) CreateBlog(ctx *gin.Context) {
	var req CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Uploading image
	// Get file extension
	fileExtension, err := util.GetFileExtensionFromBase64(req.Image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate file name
	fileName := uuid.New().String()

	// Convert image.ImageUrl (Base64) to image file
	fileBase64 := util.GetBase64Data(req.Image)

	// Convert image.ImageUrl (Base64) to image file
	image_location := fmt.Sprintf("./image/%s%s", fileName, fileExtension)
	util.SaveBase64ToFile(fileBase64, image_location)
	Image_url_result, err := util.UploadFileToMinio(ctx, image_location)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete image file
	err = os.Remove(image_location)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateBlogParams{
		Title:   req.Title,
		Content: req.Content,
		Image:   *Image_url_result,
		Url:     req.URL,
	}

	// Insert Blog
	blog, err := server.store.CreateBlog(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert Blog Tag
	for _, tag := range req.BlogTags {
		if tag.Deleted {
			continue
		}

		err := server.store.CreateBlogTag(ctx, db.CreateBlogTagParams{
			BlogID: blog.ID,
			TagID:  tag.TagId,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// Get BlogTags by Blog ID
	blogTags, err := server.store.GetBlogTagByBlogId(ctx, blog.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data: CreateBlogByIdResponse{
			Blog:     blog,
			BlogTags: blogTags,
		},
	})
}

type UpdateBlogRequest struct {
	ID       int64            `json:"id" binding:"required,min=1"`
	Title    string           `json:"title" binding:"required"`
	Content  string           `json:"content" binding:"required"`
	Image    string           `json:"image" binding:"required"`
	URL      string           `json:"url" binding:"required"`
	BlogTags []BlogTagRequest `json:"blog_tags"`
}

type UpdateBlogResponse struct {
	db.GetBlogByIdRow
	BlogTags []db.GetBlogTagByBlogIdRow `json:"blog_tags"`
}

// UpdateBlog godoc
//
//	@Summary		Update Blog
//	@Description	Update Blog
//	@Tags			Blog
//	@Accept			json
//	@Produce		json
//	@Param			input	body		UpdateBlogRequest	true	"Update information"
//	@Success		200		{object}	jsonResponse
//	@Router			/api/blog [put]
//	@Security		BearerAuth
func (server *Server) UpdateBlog(ctx *gin.Context) {
	var req UpdateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateBlogParams{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Image:   req.Image,
		Url:     req.URL,
	}

	// Check Image is link or base64
	if !util.IsValidURL(req.Image) {
		// Uploading image
		// Get file extension
		fileExtension, err := util.GetFileExtensionFromBase64(req.Image)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Generate file name
		fileName := uuid.New().String()

		// Convert image.ImageUrl (Base64) to image file
		fileBase64 := util.GetBase64Data(req.Image)

		// Convert image.ImageUrl (Base64) to image file
		image_location := fmt.Sprintf("./image/%s%s", fileName, fileExtension)
		util.SaveBase64ToFile(fileBase64, image_location)
		Image_url_result, err := util.UploadFileToMinio(ctx, image_location)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Delete image file
		err = os.Remove(image_location)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg.Image = *Image_url_result
	}

	// Update Blog
	err := server.store.UpdateBlog(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, bt := range req.BlogTags {

		argGetBlogTagByBlogIdAndTagId := db.GetBlogTagByBlogIdAndTagIdParams{
			BlogID: bt.BlogID,
			TagID:  bt.TagId,
		}

		// Check exists blog_tag
		blogTagExists, _ := server.store.GetBlogTagByBlogIdAndTagId(ctx, argGetBlogTagByBlogIdAndTagId)
		if blogTagExists.ID == 0 && !bt.Deleted {
			// Insert Blog Tag
			err := server.store.CreateBlogTag(ctx, db.CreateBlogTagParams{
				BlogID: req.ID,
				TagID:  bt.TagId,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else if blogTagExists.ID != 0 && bt.Deleted {
			// Delete Blog Tag
			err := server.store.DeleteBlogTag(ctx, db.DeleteBlogTagParams{
				BlogID: req.ID,
				TagID:  bt.TagId,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	// Get Blog by Id
	blog, err := server.store.GetBlogById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get Blog Tags
	blogTags, err := server.store.GetBlogTagByBlogId(ctx, blog.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data: UpdateBlogResponse{
			GetBlogByIdRow: blog,
			BlogTags:       blogTags,
		},
	})
}

type DeleteBlogRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// DeleteBlog godoc
//
//	@Summary		Delete Blog
//	@Description	Delete Blog
//	@Tags			Blog
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Blog ID"
//	@Success		200	{object}	jsonResponse
//	@Router			/api/blog/{id} [delete]
//	@Security		BearerAuth
func (server *Server) DeleteBlog(ctx *gin.Context) {
	var req DeleteBlogRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Delete Blog
	err := server.store.DeleteBlog(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete Blog Tag by Blog ID
	err = server.store.DeleteBlogTagByBlogId(ctx, req.ID)
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
