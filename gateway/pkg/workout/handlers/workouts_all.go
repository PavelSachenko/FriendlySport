package handlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/utils"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"net/http"
)

type workoutsFiltering struct {
	Title         *string `form:"title" json:"title"`
	IsDone        *bool   `form:"is_done" json:"is_done"`
	AppointedTime *string `form:"appointed_time" json:"appointed_time"`
	Sort          *string `form:"sort" json:"sort"`
	Limit         uint64  `form:"limit,default=10" json:"limit"`
	Offset        uint64  `form:"offset,default=0" json:"offset"`
}

func WorkoutsAll(ctx *gin.Context, c workout.WorkoutServiceClient) {
	var workoutsFiltering workoutsFiltering
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}
	err := ctx.ShouldBindQuery(&workoutsFiltering)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	query, _ := json.Marshal(workoutsFiltering)
	res, err := c.All(context.Background(), &workout.WorkoutFilteringRequest{
		Query:  query,
		UserId: userId,
	})
	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}
	ctx.JSON(http.StatusOK, res)
}
