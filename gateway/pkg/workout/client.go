package workout

import (
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/workout/pb/workout"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	Workout  workout.WorkoutServiceClient
	Exercise workout.ExerciseServiceClient
}

func InitServiceClient(cfg config.Config) ServiceClient {
	log.Printf("Initial workout service")
	cc, err := grpc.Dial(cfg.WorkoutServiceUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can not connect to user service: %v", err)
	}

	return ServiceClient{
		Workout:  workout.NewWorkoutServiceClient(cc),
		Exercise: workout.NewExerciseServiceClient(cc),
	}
}
