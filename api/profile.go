package api

import (
	db "blog-go-api/db/sqlc"
	"blog-go-api/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Get Profile
// @Description	Get Profile
// @Tags			Profile
// @Accept			json
// @Produce		json
// @Success		200	{object}	userResponse
// @Router			/api/profile [get]
// @Security		BearerAuth
func (server *Server) GetProfile(ctx *gin.Context) {
	userId, err := getUserFromContext(*server, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	profile, err := server.store.GetUser(ctx, *userId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Get profile successfully",
		Data:    profile,
	}

	ctx.JSON(http.StatusOK, payload)
}

// @Summary		Update Profile
// @Description	Update Profile
// @Tags			Profile
// @Accept			json
// @Produce		json
// @Param			input	body		updateUserRequest	true	"Update information"
// @Success		200		{object}	userResponse
// @Router			/api/profile [put]
// @Security		BearerAuth
func (server *Server) UpdateProfile(ctx *gin.Context) {
	userId, err := getUserFromContext(*server, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req db.UpdateUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.UserID = *userId

	// Update User
	_, err = server.store.UpdateUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get Updated User
	user, err := server.store.GetUser(ctx, *userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Response
	payload := jsonResponse{
		Error:   false,
		Message: "Update profile successfully",
		Data:    user,
	}

	ctx.JSON(http.StatusOK, payload)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" minLength:"8"`
	NewPassword string `json:"new_password" binding:"required" minLength:"8"`
}

// @Summary		Change Password
// @Description	Change Password
// @Tags			Profile
// @Accept			json
// @Produce		json
// @Param			input	body		ChangePasswordRequest	true	"Change Password"
// @Success		200		{object}	interface{}
// @Router			/api/profile/change_password [put]
// @Security		BearerAuth
func (server *Server) ChangePassword(ctx *gin.Context) {
	userId, err := getUserFromContext(*server, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get HashedPassword
	hashedPassword, err := server.store.GetUserHashedPassword(ctx, *userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check Password
	err = util.CheckPassword(req.OldPassword, hashedPassword)
	if err != nil {
		payload := jsonResponse{
			Error:   true,
			Message: "Old password is incorrect",
			Data:    nil,
		}
		ctx.JSON(http.StatusBadRequest, payload)
		return
	}

	// Hash New Password
	hashedPassword, err = util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Update Password
	var updateUserPasswordParams db.UpdateUserPasswordParams
	updateUserPasswordParams.ID = *userId
	updateUserPasswordParams.HashedPassword = hashedPassword

	err = server.store.UpdateUserPassword(ctx, updateUserPasswordParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Change password successfully",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, payload)
}
