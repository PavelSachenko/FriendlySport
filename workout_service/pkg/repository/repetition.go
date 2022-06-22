package repository

import (
	"context"
	"github.com/pavel/workout_service/pkg/model"
)

type Repetition interface {
	Create(repetition model.Repetition) (error error, createdExercise *model.Repetition)
	Delete(ctx context.Context) error
	Update(ctx context.Context, update model.ExerciseUpdate) (error, *model.Repetition)
}
