package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/user/pb/auth"
	"net/http"
)

type registerRequest struct {
	Email          string `json:"email"   validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6,max=60"`
	RepeatPassword string `json:"repeat_password" validate:"required,min=6,max=60,eqfield=Password"`
	RoleId         uint32 `json:"role_id" validate:"required"`
}

func Register(ctx *gin.Context, c auth.AuthServiceClient) {
	var registerRequest registerRequest
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &registerRequest)

	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}
	res, err := c.Register(context.Background(), &auth.RegisterRequest{
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
		RoleId:   registerRequest.RoleId,
	})

	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
