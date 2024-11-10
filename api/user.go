package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "blog-go-api/db/sqlc"
	"blog-go-api/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createUserRequest struct {
	Username  string `json:"username" binding:"required,alphanum"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,min=10"`
}

type userResponse struct {
	Code              string    `json:"code"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Description       string    `json:"description"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Code:              user.Code,
		Username:          user.Username,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		Phone:             user.Phone,
		Description:       user.Description.String,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// @Summary		Sign up
// @Description	Create a new user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			input	body		createUserRequest	true	"User information"
// @Success		200		{object}	jsonResponse
// @Router			/api/signup [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Generate Hash Password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get Count User
	countUser, err := server.store.CountUserForGenerateCode(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate User Code
	code := fmt.Sprintf("U%06d", countUser+1)

	arg := db.CreateUserParams{
		Code:           code,
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.FirstName,
		Email:          req.Email,
		Phone:          req.Phone,
		Description: pgtype.Text{
			String: "",
			Valid:  false,
		},
	}

	// Create User
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create User Role
	argCreateUserRole := db.CreateUserRoleParams{
		UserID: user.ID,
		RoleID: 1,
	}

	err = server.store.CreateUserRole(ctx, argCreateUserRole)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result := jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, result)
}

type updateUserRequest struct {
	ID          int64   `json:"id" binding:"required,min=1"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email" binding:"email"`
	Phone       string  `json:"phone"`
	Description string  `json:"description"`
	Roles       []int64 `json:"roles"`
}

type updateResponse struct {
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Description       string    `json:"description"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// @Summary		Update user
// @Description	Update a user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			input	body		updateUserRequest	true	"User information"
// @Success		200		{object}	jsonResponse
// @Router			/api/users [put]
// @Security		BearerAuth
func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		UserID: req.ID,
		FirstName: pgtype.Text{
			String: req.FirstName,
			Valid:  req.FirstName != "",
		},
		LastName: pgtype.Text{
			String: req.LastName,
			Valid:  req.LastName != "",
		},
		Email: pgtype.Text{
			String: req.Email,
			Valid:  req.Email != "",
		},
		Phone: pgtype.Text{
			String: req.Phone,
			Valid:  req.Phone != "",
		},
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Phone != "",
		},
	}

	// if req.Password != "" {
	// 	hashedPassword, err := util.HashPassword(req.Password)
	// 	if err != nil {
	// 		err := fmt.Errorf("failed to hash password: %s", err)
	// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 		return
	// 	}

	// 	arg.HashedPassword = pgtype.Text{
	// 		String: hashedPassword,
	// 		Valid:  true,
	// 	}

	// 	arg.PasswordChangedAt = pgtype.Timestamptz{
	// 		Time:  time.Now(),
	// 		Valid: true,
	// 	}
	// }

	// Update User
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("user not found: %s", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		err = fmt.Errorf("failed to update user: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Update User Role
	// Delete All User Role
	err = server.store.DeleteUserRoleByUserId(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create User Role
	for _, roleID := range req.Roles {
		argCreateUserRole := db.CreateUserRoleParams{
			UserID: user.ID,
			RoleID: roleID,
		}

		err := server.store.CreateUserRole(ctx, argCreateUserRole)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	result := jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, result)
}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

type listUserResponse struct {
	db.ListUsersRow
	Roles []db.GetRoleByUserIdRow `json:"roles"`
}

// @Summary		List users
// @Description	Get a list of users
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			page_id		query		int	true	"Page ID"
// @Param			page_size	query		int	true	"Page Size"
// @Success		200			{object}	jsonResponseWithPaginate
// @Security		BearerAuth
// @Router			/api/users [get]
func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	var users_response []listUserResponse

	user, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	countUser, err := server.store.CountUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get Role By User ID
	for _, u := range user {
		roles, err := server.store.GetRoleByUserId(ctx, u.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newUserResponse := listUserResponse{
			ListUsersRow: u,
			Roles:        roles,
		}

		users_response = append(users_response, newUserResponse)
	}

	payload := jsonResponseWithPaginate{
		jsonResponse: jsonResponse{
			Error:   false,
			Message: "successfully",
			Data:    users_response,
		},
		Total: countUser,
	}

	ctx.JSON(http.StatusOK, payload)
}

type getUser struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserResponse struct {
	db.GetUserRow
	Roles []db.GetRoleByUserIdRow `json:"roles"`
}

// @Summary		Get user
// @Description	Get a user by ID
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	jsonResponse
// @Router			/api/users/{id} [get]
// @Security		BearerAuth
func (server *Server) getUser(ctx *gin.Context) {
	var req getUser
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	User, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get role
	roles, err := server.store.GetRoleByUserId(ctx, User.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user_response := getUserResponse{
		GetUserRow: User,
		Roles:      roles,
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// if account.Owner != authPayload.Username {
	// 	err := errors.New("account doesn't belong to the authenticated user")
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }

	payload := jsonResponse{
		Error:   false,
		Message: "Get user successfully",
		Data:    user_response,
	}

	ctx.JSON(http.StatusOK, payload)
}

type deleteUser struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary		Delete user
// @Description	Delete a user by ID
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	jsonResponse
// @Router			/api/users/{id} [delete]
// @Security		BearerAuth
func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUser
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}

	// Delete User
	err = server.store.DeleteUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete User Role
	err = server.store.DeleteUserRoleByUserId(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// if account.Owner != authPayload.Username {
	// 	err := errors.New("account doesn't belong to the authenticated user")
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }

	payload := jsonResponse{
		Error:   false,
		Message: "Delete user successfully",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, payload)
}
