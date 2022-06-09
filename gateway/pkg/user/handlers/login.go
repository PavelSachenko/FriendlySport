package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/user/pb/auth"
	"net/http"
	"time"
)

type loginRequest struct {
	Password string `json:"password"  validate:"required,min=5,max=500"`
	Email    string `json:"email"   validate:"required,min=5,max=500"`
}

func Login(ctx *gin.Context, c auth.AuthServiceClient) {
	var loginRequest loginRequest
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &loginRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}
	res, err := c.Login(context.Background(), &auth.LoginRequest{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Duration(res.RefreshTokenExpire)),
	})
	ctx.JSON(http.StatusOK, res)
}
