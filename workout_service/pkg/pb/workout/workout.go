package workout

import (
	"context"
	"encoding/json"
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
	var t time.Time
	t = time.Unix(request.AppointedTime, 0)
	err, workout := s.workout.AddList(&model.Workout{
		UserId:        request.UserId,
		Title:         request.Title,
		Description:   request.Description,
		AppointedTime: &t,
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
	err := s.workout.DeleteList(request.Id, request.UserId)
	if err != nil {
		return &DeleteResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}
	return &DeleteResponse{
		Status: http.StatusNoContent,
	}, nil
}

func (s Server) Update(ctx context.Context, request *UpdateRequest) (*UpdateResponse, error) {
	workoutUpdate := model.WorkoutUpdate{UserId: request.UserId, Id: request.Id, UpdatedAt: time.Now()}
	json.Unmarshal(request.Query, &workoutUpdate)
	err, res := s.workout.UpdateList(workoutUpdate)

	if err != nil {
		return &UpdateResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}, nil
	}
	return &UpdateResponse{
		Status: http.StatusOK,
		Workout: &Workout{
			Id:            res.ID,
			UserId:        res.UserId,
			Title:         res.Title,
			Description:   res.Description,
			IsDone:        res.IsDone,
			AppointedTime: res.AppointedTime.Unix(),
			CreatedAt:     res.CreatedAt.Unix(),
			UpdatedAt:     res.UpdatedAt.Unix(),
		},
	}, nil
}

func (s Server) All(ctx context.Context, request *WorkoutFilteringRequest) (*WorkoutFilteringResponse, error) {
	var workoutsFiltering model.WorkoutsFiltering
	json.Unmarshal(request.Query, &workoutsFiltering)
	err, res := s.workout.GetAll(request.UserId, workoutsFiltering)
	if err != nil {
		return &WorkoutFilteringResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}, nil
	}

	var workoutsList []*Workout
	for _, workout := range res {
		var appointedTime int64
		if workout.AppointedTime != nil {
			appointedTime = workout.AppointedTime.Unix()
		}
		jsonExercises, _ := json.Marshal(workout.Exercises)
		var exercises []*Exercises
		err = json.Unmarshal(jsonExercises, &exercises)
		workoutsList = append(workoutsList, &Workout{
			Id:            workout.ID,
			Title:         workout.Title,
			Description:   workout.Description,
			UserId:        workout.UserId,
			IsDone:        workout.IsDone,
			AppointedTime: appointedTime,
			CreatedAt:     workout.CreatedAt.Unix(),
			UpdatedAt:     workout.UpdatedAt.Unix(),
			Exercise:      exercises,
		})
	}

	return &WorkoutFilteringResponse{
		Status:  http.StatusOK,
		Workout: workoutsList,
	}, nil
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
