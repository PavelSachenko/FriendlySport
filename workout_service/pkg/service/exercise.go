package service

import (
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/pkg/repository"
)

type Exercise interface {
	AddToWorkoutList(userId uint64, exercise model.Exercise) (error, *model.Exercise)
	UpdateInWorkoutList(exercise model.ExerciseUpdate) (error, *model.Exercise)
	DeleteFromWorkoutList(id, userId, workoutId uint64) error
	GetRecommendationTitle(typingTitle string) (error, []*model.ExerciseRecommendation)
}

type ExerciseService struct {
	repo repository.Exercise
}

func InitExerciseService(repo repository.Exercise) *ExerciseService {
	return &ExerciseService{
		repo: repo,
	}
}

func (e ExerciseService) AddToWorkoutList(userId uint64, exercise model.Exercise) (error, *model.Exercise) {
	return e.repo.Create(userId, exercise)
}

func (e ExerciseService) UpdateInWorkoutList(exercise model.ExerciseUpdate) (error, *model.Exercise) {
	return e.repo.Update(exercise)
}

func (e ExerciseService) DeleteFromWorkoutList(id, userId, workoutId uint64) error {
	return e.repo.Delete(id, userId, workoutId)
}

func (e ExerciseService) GetRecommendationTitle(typingTitle string) (error, []*model.ExerciseRecommendation) {
	return e.repo.GetRecommendation(typingTitle)
}
