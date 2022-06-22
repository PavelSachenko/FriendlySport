package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/user/pb/auth"
	"net/http"
	"time"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func RefreshAuthToken(ctx *gin.Context, c auth.AuthServiceClient) {
	var refreshRequest refreshRequest
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &refreshRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}
	res, err := c.RefreshAuthToken(context.Background(), &auth.RefreshTokenRequest{
		RefreshToken: refreshRequest.RefreshToken,
	})

	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
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
