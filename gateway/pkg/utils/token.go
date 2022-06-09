package utils

import (
	"github.com/gin-gonic/gin"
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
