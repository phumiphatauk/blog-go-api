package api

import (
	db "blog-go-api/db/sqlc"
	mail "blog-go-api/mail"
	"blog-go-api/util"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
	Permissions           []string     `json:"permissions"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get User
	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check Password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Create Access Token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create Refresh Token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create Session
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get Permissions
	permissions, err := server.store.GetPermissionByUserId(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Response
	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
		Permissions:           permissions,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (server *Server) forgotPassword(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	config, err := util.LoadConfig(".")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get User by Email
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create Reset Password
	token := uuid.New().String()
	_, err = server.store.CreateResetPassword(ctx, db.CreateResetPasswordParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24), Valid: true},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sender := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Confirm reset password"
	content := fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Password Reset</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
			<div style="max-width: 600px; margin: auto; background-color: #ffffff; padding: 20px; border-radius: 10px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">

				<div style="text-align: center; margin-bottom: 20px;">
					<img src="https://eslens.online/images/logo-dark.png" alt="Company Logo" style="max-width: 150px;">
				</div>

				<h2 style="color: #333333;">Reset Your Password</h2>
				<p>Dear %s,</p>
				<p>We received a request to reset the password for your account. If you did not request this, please ignore this email.</p>
				<p>To reset your password, please click the button below:</p>
				<p style="text-align: center;">
					<a href="https://eslens.online/reset-password/%s" target="_blank" style="display: inline-block; background-color: #17A34A; color: #ffffff; padding: 10px 20px; text-decoration: none; border-radius: 5px; font-size: 16px;cursor:pointer">
						Reset Password
					</a>
				</p>
				<p>This link will expire in 24 hours. If it expires, you will need to request another password reset.</p>
				<p>If you have any questions or need assistance, feel free to contact our support team.</p>
				<p>Best regards,</p>
				<p>Automated System<br>
				Eslens.Co.,Ltd.<br></p>
			</div>
		</body>
		</html>`, user.FirstName, token)

	to := []string{req.Email}

	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(nil))
}

type resetPasswordDetailRequest struct {
	Token string `form:"token" binding:"required"`
}

type resetPasswordDetailResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (server *Server) resetPasswordDetail(ctx *gin.Context) {
	var req resetPasswordDetailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get Reset Password by Token
	resetPassword, err := server.store.GetResetPasswordByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check Token Expired
	if resetPassword.ExpiresAt.Time.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("token expired")))
		return
	}

	// Get User by ID
	user, err := server.store.GetUser(ctx, resetPassword.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(resetPasswordDetailResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}))
}

type resetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get Reset Password by Token
	resetPassword, err := server.store.GetResetPasswordByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check Token Expired
	if resetPassword.ExpiresAt.Time.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("token expired")))
		return
	}

	// Get User by ID
	user, err := server.store.GetUser(ctx, resetPassword.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Hash Password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Update Password
	err = server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:             user.ID,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Set Used for Reset Password
	err = server.store.UseResetPassword(ctx, req.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(nil))
}
