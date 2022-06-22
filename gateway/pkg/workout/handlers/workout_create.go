package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/utils"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"net/http"
)

type workoutCreateRequest struct {
	Title         string `json:"title"   validate:"required,min=3,max=255"`
	Description   string `json:"description"   validate:"omitempty,min=5,max=500"`
	AppointedTime int64  `json:"appointed_time"   validate:"omitempty,number"`
}

func WorkoutCreate(ctx *gin.Context, c workout.WorkoutServiceClient) {
	var workoutCreateRequest workoutCreateRequest
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &workoutCreateRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}

	res, err := c.Create(context.Background(), &workout.CreateWorkoutRequest{
		UserId:        userId,
		Title:         workoutCreateRequest.Title,
		Description:   workoutCreateRequest.Description,
		AppointedTime: workoutCreateRequest.AppointedTime,
	})

	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
