package pb

import (
	"context"
	"github.com/pavel/workout_service/pkg/adapter"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/pb/workout"
	"github.com/pavel/workout_service/pkg/service"
	"net/http"
)

type GRPCExerciseService struct {
	exercise service.Exercise
	logger   *logger.Logger
	workout.ExerciseServiceServer
}

func InitGRPCExerciseService(exercise service.Exercise, logger *logger.Logger) *GRPCExerciseService {
	return &GRPCExerciseService{
		exercise: exercise,
		logger:   logger,
	}
}

func (s GRPCExerciseService) Create(ctx context.Context, request *workout.CreateExerciseRequest) (*workout.CreateExerciseResponse, error) {
	err, e := s.exercise.AddToWorkoutList(request.UserId, adapter.GRPCToExercise(request, s.logger))
	if err != nil {
		return &workout.CreateExerciseResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}
	return &workout.CreateExerciseResponse{
		Status:   http.StatusCreated,
		Exercise: adapter.ExerciseToGRPC(*e, s.logger),
	}, nil
}

func (s GRPCExerciseService) Update(ctx context.Context, request *workout.UpdateExerciseRequest) (*workout.UpdateExerciseResponse, error) {
	errRes, e := adapter.GRPCExerciseUpdateToExerciseUpdate(request, s.logger)
	if errRes != nil {
		return errRes, nil
	}
	err, res := s.exercise.UpdateInWorkoutList(*e)
	if err != nil {
		return &workout.UpdateExerciseResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &workout.UpdateExerciseResponse{
		Status:   http.StatusCreated,
		Exercise: adapter.ExerciseToGRPC(*res, s.logger),
	}, nil
}

func (s GRPCExerciseService) Delete(ctx context.Context, request *workout.DeleteExerciseRequest) (*workout.DeleteExerciseResponse, error) {
	s.logger.Infof("exercise:delete id: %d", request.Id)
	err := s.exercise.DeleteFromWorkoutList(request.Id, request.UserId, request.WorkoutId)
	if err != nil {
		return &workout.DeleteExerciseResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}
	return &workout.DeleteExerciseResponse{
		Status: http.StatusNoContent,
	}, nil
}

func (s GRPCExerciseService) ExerciseTitleRecommendation(ctx context.Context, request *workout.ExerciseTitleRecommendationRequest) (*workout.ExerciseTitleRecommendationResponse, error) {
	err, recommendationList := s.exercise.GetRecommendationTitle(request.TypingTitle)
	if err != nil {
		return &workout.ExerciseTitleRecommendationResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &workout.ExerciseTitleRecommendationResponse{
		Status:             http.StatusOK,
		RecommendationList: adapter.ExerciseRecommendationTitleToGRPC(recommendationList, s.logger),
	}, nil
}
