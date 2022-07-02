package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pavel/gateway/pkg/user/pb/role"
	"net/http"
)

func Roles(ctx *gin.Context, c role.RoleServiceClient) {
	res, err := c.All(context.Background(), &empty.Empty{})
	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}
	ctx.JSON(http.StatusOK, res)
}
