package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"net/http"
	"strings"
)

func GetBearerToken(ctx *gin.Context) string {
	authorization := ctx.Request.Header.Get("authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return token[1]
}

func GetUserIdFromContext(ctx *gin.Context) (*errors.IError, uint64) {
	userId, ok := ctx.Get("userId")
	if ok == false {
		return &errors.IError{Field: "Bearer Token", Value: "wrong", Tag: "token"}, 0
	}
	return nil, userId.(uint64)
}
