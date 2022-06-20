package adapter

import (
	"encoding/json"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/pkg/pb/workout"
	"net/http"
	"time"
)

func WorkoutListToGRPC(workouts []*model.Workout, logger logger.Logger) (*workout.WorkoutFilteringResponse, []*workout.Workout) {
	workoutsJson, err := json.Marshal(workouts)
	if err != nil {
		logger.Error(err)
		if err != nil {
			logger.Error(err)
			return &workout.WorkoutFilteringResponse{
				Error:  err.Error(),
				Status: http.StatusUnprocessableEntity,
			}, nil
		}
	}
	var workoutsList []*workout.Workout
	err = json.Unmarshal(workoutsJson, &workoutsList)
	if err != nil {
		logger.Error(err)
		if err != nil {
			logger.Error(err)
			return &workout.WorkoutFilteringResponse{
				Error:  err.Error(),
				Status: http.StatusUnprocessableEntity,
			}, nil
		}
	}

	logger.Infof("workouts:get-all response: %v", workoutsList)

	return nil, workoutsList
}

func GRPCToWorkoutList(request *workout.WorkoutFilteringRequest, logger logger.Logger) (*workout.WorkoutFilteringResponse, *model.WorkoutsFiltering) {
	var workoutsFiltering model.WorkoutsFiltering
	err := json.Unmarshal(request.Query, &workoutsFiltering)
	if err != nil {
		logger.Error(err)
		return &workout.WorkoutFilteringResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		}, nil
	}

	logger.Infof("workouts:get-all query %v", workoutsFiltering)
	return nil, &workoutsFiltering
}

func GRPCWorkoutCreateToWorkout(request *workout.CreateRequest, logger *logger.Logger) *model.Workout {
	var t time.Time
	t = time.Unix(request.AppointedTime, 0)
	w := &model.Workout{
		UserId:        request.UserId,
		Title:         request.Title,
		Description:   request.Description,
		AppointedTime: &t,
	}

	logger.Infof("workout:create request %v", w)
	return w
}

func GRPCWorkoutUpdateToWorkoutUpdate(request *workout.UpdateRequest, logger *logger.Logger) (*workout.UpdateResponse, *model.WorkoutUpdate) {
	workoutUpdate := model.WorkoutUpdate{UserId: request.UserId, Id: request.Id, UpdatedAt: time.Now()}
	err := json.Unmarshal(request.Query, &workoutUpdate)
	if err != nil {
		logger.Error(err)
		return &workout.UpdateResponse{
			Error:  err.Error(),
			Status: http.StatusUnprocessableEntity,
		}, nil
	}

	logger.Infof("workout:update request %v", workoutUpdate)
	return nil, &workoutUpdate
}

func WorkoutToGRPC(model *model.Workout, logger *logger.Logger) *workout.Workout {
	w := &workout.Workout{
		Id:            model.ID,
		UserId:        model.UserId,
		Title:         model.Title,
		Description:   model.Description,
		IsDone:        model.IsDone,
		AppointedTime: model.AppointedTime.Unix(),
		CreatedAt:     model.CreatedAt.Unix(),
		UpdatedAt:     model.UpdatedAt.Unix(),
	}

	logger.Infof("workout response %v", w)
	return w
}

func WorkoutRecommendationTitleToGRPC(recommendationList []*model.WorkoutRecommendation, logger *logger.Logger) []*workout.WorkoutTitleRecommendation {
	var recommendations []*workout.WorkoutTitleRecommendation
	for _, recommendation := range recommendationList {
		recommendations = append(recommendations, &workout.WorkoutTitleRecommendation{Title: recommendation.Title})
	}

	return recommendations
}
