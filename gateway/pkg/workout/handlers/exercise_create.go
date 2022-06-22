package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/utils"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"net/http"
	"strconv"
)

type exerciseCreateRequest struct {
	Title       string `json:"title"   validate:"required,min=3,max=255"`
	Description string `json:"description"   validate:"omitempty,min=5,max=500"`
}

func ExerciseCreate(ctx *gin.Context, c workout.ExerciseServiceClient) {
	var exerciseCreateRequest exerciseCreateRequest
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &exerciseCreateRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}
	workoutId, err := strconv.ParseUint(ctx.Param("workout_id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	res, err := c.Create(context.Background(), &workout.CreateExerciseRequest{
		UserId:      userId,
		WorkoutId:   workoutId,
		Title:       exerciseCreateRequest.Title,
		Description: exerciseCreateRequest.Description,
	})

	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
