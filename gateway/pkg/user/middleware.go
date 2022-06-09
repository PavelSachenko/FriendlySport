package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/user/pb/auth"
	"github.com/pavel/gateway/pkg/utils"
	"net/http"
)

type AuthMiddleware struct {
	svc ServiceClient
}

func (c *AuthMiddleware) AuthRequired(ctx *gin.Context) {
	res, err := c.svc.Auth.CheckAuthToken(context.Background(), &auth.CheckTokenRequest{
		Token: utils.GetBearerToken(ctx),
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.IError{Field: "Bearer Token", Value: "wrong", Tag: "token"})
		return
	}

	ctx.Set("userId", res.UserId)

	ctx.Next()
}
