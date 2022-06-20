package pb

import (
	"context"
	"github.com/pavel/workout_service/pkg/adapter"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/pb/workout"
	"github.com/pavel/workout_service/pkg/service"
	"net/http"
)

type Server struct {
	workout service.Workout
	logger  *logger.Logger
	workout.WorkoutServiceServer
}

func InitGRPCWorkoutServer(workout service.Workout, logger *logger.Logger) *Server {
	return &Server{
		workout: workout,
		logger:  logger,
	}
}

func (s Server) Create(context context.Context, request *workout.CreateRequest) (*workout.CreateResponse, error) {

	w := adapter.GRPCWorkoutCreateToWorkout(request, s.logger)
	err, res := s.workout.AddList(w)
	if err != nil {
		return &workout.CreateResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &workout.CreateResponse{
		Status:  http.StatusCreated,
		Workout: adapter.WorkoutToGRPC(res, s.logger),
	}, nil
}

func (s Server) Delete(context context.Context, request *workout.DeleteRequest) (*workout.DeleteResponse, error) {
	err := s.workout.DeleteList(request.Id, request.UserId)
	if err != nil {
		return &workout.DeleteResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}
	return &workout.DeleteResponse{
		Status: http.StatusNoContent,
	}, nil
}

func (s Server) Update(context context.Context, request *workout.UpdateRequest) (*workout.UpdateResponse, error) {
	wError, workoutUpdate := adapter.GRPCWorkoutUpdateToWorkoutUpdate(request, s.logger)
	if wError != nil {
		return wError, nil
	}
	err, res := s.workout.UpdateList(*workoutUpdate)

	if err != nil {
		return &workout.UpdateResponse{
			Error:  err.Error(),
			Status: http.StatusUnprocessableEntity,
		}, nil
	}
	return &workout.UpdateResponse{
		Status:  http.StatusOK,
		Workout: adapter.WorkoutToGRPC(res, s.logger),
	}, nil
}

func (s Server) All(context context.Context, request *workout.WorkoutFilteringRequest) (*workout.WorkoutFilteringResponse, error) {
	resError, workoutsFiltering := adapter.GRPCToWorkoutList(request, *s.logger)
	if resError != nil {
		return resError, nil
	}

	err, res := s.workout.GetAll(request.UserId, *workoutsFiltering)
	if err != nil {
		s.logger.Error(err)
		return &workout.WorkoutFilteringResponse{
			Error:  err.Error(),
			Status: http.StatusNotFound,
		}, nil
	}
	resError, workoutsLists := adapter.WorkoutListToGRPC(res, *s.logger)

	return &workout.WorkoutFilteringResponse{
		Status:  http.StatusOK,
		Workout: workoutsLists,
	}, nil
}

func (s Server) WorkoutTitleRecommendation(context context.Context, request *workout.WorkoutTitleRecommendationRequest) (*workout.WorkoutTitleRecommendationResponse, error) {
	err, recommendationList := s.workout.RecommendationTitles(request.TypingTitle)
	if err != nil {
		return &workout.WorkoutTitleRecommendationResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &workout.WorkoutTitleRecommendationResponse{
		Status:             http.StatusOK,
		RecommendationList: adapter.WorkoutRecommendationTitleToGRPC(recommendationList, s.logger),
	}, nil
}
