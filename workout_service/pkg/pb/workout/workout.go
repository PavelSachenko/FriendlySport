package workout

import (
	"context"
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/pkg/service"
	"net/http"
	"time"
)

type Server struct {
	workout service.Workout
}

func InitGRPCWorkoutServer(workout service.Workout) *Server {
	return &Server{
		workout: workout,
	}
}

func (s Server) Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error) {
	err, workout := s.workout.AddList(&model.Workout{
		UserId:        request.UserId,
		Title:         request.Title,
		Description:   request.Description,
		AppointedTime: time.Unix(request.AppointedTime, 0),
	})
	if err != nil {
		return &CreateResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &CreateResponse{
		Status: http.StatusCreated,
		Workout: &Workout{
			Id:            workout.ID,
			UserId:        workout.UserId,
			Title:         workout.Title,
			Description:   workout.Description,
			IsDone:        workout.IsDone,
			AppointedTime: workout.AppointedTime.Unix(),
			CreatedAt:     workout.CreatedAt.Unix(),
			UpdatedAt:     workout.UpdatedAt.Unix(),
		},
	}, nil
}

func (s Server) Delete(ctx context.Context, request *DeleteRequest) (*DeleteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Update(ctx context.Context, request *UpdateRequest) (*UpdateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) One(ctx context.Context, request *OneRequest) (*OneResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) All(ctx context.Context, request *AllRequest) (*AllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) mustEmbedUnimplementedWorkoutServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s Server) WorkoutTitleRecommendation(ctx context.Context, request *WorkoutTitleRecommendationRequest) (*WorkoutTitleRecommendationResponse, error) {
	err, recommendationList := s.workout.RecommendationTitles(request.TypingTitle)
	if err != nil {
		return &WorkoutTitleRecommendationResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	var recommendations []*WorkoutTitleRecommendation
	for _, recommendation := range recommendationList {
		recommendations = append(recommendations, &WorkoutTitleRecommendation{Title: recommendation.Title})
	}

	return &WorkoutTitleRecommendationResponse{
		Status:             http.StatusOK,
		RecommendationList: recommendations,
	}, nil
}
