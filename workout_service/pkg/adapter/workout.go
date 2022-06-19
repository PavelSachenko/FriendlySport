package adapter

import (
	"encoding/json"
	"github.com/pavel/workout_service/pkg/errors"
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/pkg/pb/workout"
)

func WorkoutListToGRPC(workouts []*model.Workout) ([]*workout.Workout, error) {
	workoutsJson, err := json.Marshal(workouts)
	if err != nil {
		return nil, errors.UnprocessableEntity
	}
	var workoutsList []*workout.Workout
	err = json.Unmarshal(workoutsJson, &workoutsList)
	if err != nil {
		return nil, errors.UnprocessableEntity
	}
	return workoutsList, nil
}

type Test interface {
	Test()
}