package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/user/pb/user"
	"github.com/pavel/gateway/pkg/utils"
	"net/http"
)

func User(ctx *gin.Context, u user.UserServiceClient) {

	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}

	res, err := u.One(context.Background(), &user.OneRequest{
		UserId: userId,
	})

	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
