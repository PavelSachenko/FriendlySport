package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/user/pb/user"
	"net/http"
)

func User(ctx *gin.Context, u user.UserServiceClient) {

	userId, ok := ctx.Get("userId")
	if ok == false {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	test := userId.(uint64)
	res, err := u.One(context.Background(), &user.OneRequest{
		UserId: test,
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
