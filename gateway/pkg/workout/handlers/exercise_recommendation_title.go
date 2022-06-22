package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/pkg/errors"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"net/http"
)

type exerciseRecommendationRequest struct {
	Title string `form:"title"   validate:"required,max=255"`
}

func ExerciseRecommendationTitle(ctx *gin.Context, c workout.ExerciseServiceClient) {
	var exerciseRecommendation exerciseRecommendationRequest
	err := ctx.ShouldBindQuery(&exerciseRecommendation)
	validate := errors.InitValidator()
	requestErrors := validate.ValidateRequest(ctx, &exerciseRecommendation)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	res, err := c.ExerciseTitleRecommendation(context.Background(), &workout.ExerciseTitleRecommendationRequest{
		TypingTitle: exerciseRecommendation.Title,
	})

	if err != nil || res.Status >= http.StatusBadRequest {
		ctx.AbortWithStatusJSON(int(res.Status), res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
