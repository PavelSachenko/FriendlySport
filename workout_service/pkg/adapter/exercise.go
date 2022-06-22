package adapter

import (
	"encoding/json"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/pkg/pb/workout"
	"net/http"
)

func GRPCToExercise(request *workout.CreateExerciseRequest, logger *logger.Logger) model.Exercise {
	e := model.Exercise{
		Title:       request.Title,
		Description: request.Description,
		WorkoutId:   request.WorkoutId,
	}
	logger.Infof("exercise:create-request %v", e)

	return e
}

func ExerciseToGRPC(exercise model.Exercise, logger *logger.Logger) *workout.Exercise {
	e := &workout.Exercise{
		Id:          exercise.ID,
		Title:       exercise.Title,
		Description: exercise.Description,
		IsDone:      exercise.IsDone,
	}
	logger.Infof("exercise:create-response %v", e)

	return e
}

func GRPCExerciseUpdateToExerciseUpdate(request *workout.UpdateExerciseRequest, logger *logger.Logger) (*workout.UpdateExerciseResponse, *model.ExerciseUpdate) {
	workoutUpdate := model.ExerciseUpdate{Id: request.Id, UserId: request.UserId, WorkoutId: request.WorkoutId}
	err := json.Unmarshal(request.Query, &workoutUpdate)
	if err != nil {
		logger.Error(err)
		return &workout.UpdateExerciseResponse{
			Error:  err.Error(),
			Status: http.StatusUnprocessableEntity,
		}, nil
	}

	logger.Infof("workout:update request %v", workoutUpdate)
	return nil, &workoutUpdate
}

func ExerciseRecommendationTitleToGRPC(recommendationList []*model.ExerciseRecommendation, logger *logger.Logger) []*workout.TitleRecommendation {
	var recommendations []*workout.TitleRecommendation
	for _, recommendation := range recommendationList {
		recommendations = append(recommendations, &workout.TitleRecommendation{Title: recommendation.Title})
	}

	logger.Infof("recommendation:exercise list %v", recommendations)
	return recommendations
}
