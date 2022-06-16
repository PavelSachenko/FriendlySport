package service

import (
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/pkg/repository"
)

type Workout interface {
	GetAll(userId uint64, filterOption model.WorkoutsFiltering) (error, []*model.Workout)
	AddList(workout *model.Workout) (error, *model.Workout)
	UpdateList(workout model.WorkoutUpdate) (error, *model.Workout)
	DeleteList(id, userId uint64) error
	RecommendationTitles(typingTitle string) (error, []*model.WorkoutRecommendation)
}

type WorkoutService struct {
	repo repository.Workout
}

func InitWorkoutService(repo repository.Workout) *WorkoutService {
	return &WorkoutService{
		repo: repo,
	}
}

func (w *WorkoutService) RecommendationTitles(typingTitle string) (error, []*model.WorkoutRecommendation) {
	return w.repo.GetRecommendation(typingTitle)
}

func (w *WorkoutService) GetAll(userId uint64, filterOption model.WorkoutsFiltering) (error, []*model.Workout) {
	return w.repo.All(userId, filterOption)
}

func (w *WorkoutService) AddList(workout *model.Workout) (error, *model.Workout) {
	return w.repo.Create(workout)
}

func (w *WorkoutService) UpdateList(workoutUpdate model.WorkoutUpdate) (error, *model.Workout) {
	return w.repo.Update(workoutUpdate)
}

func (w *WorkoutService) DeleteList(id, userId uint64) error {
	return w.repo.Delete(id, userId)
}
