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

type exerciseUpdate struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=255"`
	Description *string `json:"description" validate:"omitempty,min=5,max=500"`
	IsDone      *bool   `json:"is_done"`
}

func ExerciseUpdate(ctx *gin.Context, c workout.ExerciseServiceClient) {
	var exerciseUpdate exerciseUpdate
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &exerciseUpdate)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	workoutId, err := strconv.ParseUint(ctx.Param("workout_id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	exerciseId, err := strconv.ParseUint(ctx.Param("exercise_id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	query, _ := json.Marshal(exerciseUpdate)
	res, err := c.Update(context.Background(), &workout.UpdateExerciseRequest{Id: exerciseId, UserId: userId, WorkoutId: workoutId, Query: query})
	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}
	ctx.JSON(http.StatusOK, res)
}
