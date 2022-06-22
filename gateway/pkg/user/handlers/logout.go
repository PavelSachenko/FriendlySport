package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/user/pb/auth"
	"github.com/pavel/gateway/pkg/utils"
	"net/http"
)

func Logout(ctx *gin.Context, c auth.AuthServiceClient) {
	token := utils.GetBearerToken(ctx)
	res, err := c.Logout(context.Background(), &auth.LogoutRequest{
		Token: token,
	})
	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   -1,
	})
	ctx.JSON(http.StatusNoContent, res)
}
