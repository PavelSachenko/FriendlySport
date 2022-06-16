package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"net/http"
)

type workoutRecommendation struct {
	Title string `form:"title"   validate:"required,max=255"`
}

func WorkoutRecommendationTitle(ctx *gin.Context, c workout.WorkoutServiceClient) {
	var workoutRecommendationRequest workoutRecommendation
	err := ctx.ShouldBindQuery(&workoutRecommendationRequest)
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &workoutRecommendationRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	res, err := c.WorkoutTitleRecommendation(context.Background(), &workout.WorkoutTitleRecommendationRequest{
		TypingTitle: workoutRecommendationRequest.Title,
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
