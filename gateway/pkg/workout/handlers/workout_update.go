package handlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/utils"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"net/http"
	"strconv"
)

type workoutUpdate struct {
	Title         *string `json:"title" validate:"omitempty,min=3,max=255"`
	Description   *string `json:"description" validate:"omitempty,min=5,max=500"`
	IsDone        *bool   `json:"is_done"`
	AppointedTime *uint64 `json:"appointed_time" validate:"omitempty,number"`
}

func WorkoutUpdate(ctx *gin.Context, c workout.WorkoutServiceClient) {
	var workoutUpdate workoutUpdate
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &workoutUpdate)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("workout_id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	query, _ := json.Marshal(workoutUpdate)
	res, err := c.Update(context.Background(), &workout.UpdateWorkoutRequest{Id: id, UserId: userId, Query: query})
	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}
	ctx.JSON(http.StatusOK, res)
}
