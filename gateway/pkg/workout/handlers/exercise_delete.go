package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/utils"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"net/http"
	"strconv"
)

func ExerciseDelete(ctx *gin.Context, c workout.ExerciseServiceClient) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
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

	res, err := c.Delete(context.Background(), &workout.DeleteExerciseRequest{Id: exerciseId, UserId: userId, WorkoutId: workoutId})
	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res.Error)
		return
	}

	ctx.JSON(http.StatusNoContent, res)
}
