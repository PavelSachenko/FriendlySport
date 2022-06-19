package pb

import (
	"context"
	workout "github.com/pavel/workout_service/pkg/pb/workout"
	"github.com/pavel/workout_service/pkg/service"
)

type Server struct {
	workout service.Workout
}

func (s Server) Create(ctx context.Context, request *workout.CreateRequest) (*workout.CreateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Delete(ctx context.Context, request *workout.DeleteRequest) (*workout.DeleteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Update(ctx context.Context, request *workout.UpdateRequest) (*workout.UpdateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) All(ctx context.Context, request *workout.WorkoutFilteringRequest) (*workout.WorkoutFilteringResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) WorkoutTitleRecommendation(ctx context.Context, request *workout.WorkoutTitleRecommendationRequest) (*workout.WorkoutTitleRecommendationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) mustEmbedUnimplementedWorkoutServiceServer() {
	//TODO implement me
	panic("implement me")
}
