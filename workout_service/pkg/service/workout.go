package service

import (
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/pkg/repository"
)

type Workout interface {
	GetOne(id uint64) (error, *model.Workout)
	GetAll()
	AddList(workout *model.Workout) (error, *model.Workout)
	UpdateList(workout *model.Workout) (error, *model.Workout)
	DeleteList(id uint64) error
}

type WorkoutService struct {
	repo repository.Workout
}

func InitWorkoutService(repo repository.Workout) *WorkoutService {
	return &WorkoutService{
		repo: repo,
	}
}

func (w *WorkoutService) GetOne(id uint64) (error, *model.Workout) {
	//TODO implement me
	panic("implement me")
}

func (w *WorkoutService) GetAll() {
	//TODO implement me
	panic("implement me")
}

func (w *WorkoutService) AddList(workout *model.Workout) (error, *model.Workout) {
	return w.repo.Create(workout)
}

func (w *WorkoutService) UpdateList(workout *model.Workout) (error, *model.Workout) {
	//TODO implement me
	panic("implement me")
}

func (w *WorkoutService) DeleteList(id uint64) error {
	//TODO implement me
	panic("implement me")
}
