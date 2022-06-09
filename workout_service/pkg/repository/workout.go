package repository

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/model"
	"time"
)

type Workout interface {
	One(id uint64) (error, *model.Workout)
	All()
	Create(workout *model.Workout) (error, *model.Workout)
	Update(workout *model.Workout) (error, *model.Workout)
	Delete(id uint64) error
}

type WorkoutRepo struct {
	*db.DB
}

func InitWorkoutRepo(db *db.DB) *WorkoutRepo {
	return &WorkoutRepo{
		DB: db,
	}
}

func (w *WorkoutRepo) One(id uint64) (error, *model.Workout) {
	//TODO implement me
	panic("implement me")
}

func (w *WorkoutRepo) All() {
	//TODO implement me
	panic("implement me")
}

func (w *WorkoutRepo) Create(workout *model.Workout) (error, *model.Workout) {
	dateNow := time.Now()
	workout.CreatedAt = dateNow
	workout.UpdatedAt = dateNow
	sql := fmt.Sprintf("INSERT INTO %s (user_id, title, description, appointed_time, created_at, updated_at) ", model.WorkoutTable)
	rows, err := w.Queryx(sql+"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		workout.UserId,
		workout.Title,
		workout.Description,
		workout.AppointedTime,
		workout.CreatedAt,
		workout.UpdatedAt,
	)
	if err != nil {
		return err, nil
	}
	for rows.Next() {
		err = rows.StructScan(&workout)
		if err != nil {
			return err, nil
		}
	}

	return nil, workout
}

func (w *WorkoutRepo) Update(workout *model.Workout) (error, *model.Workout) {
	//TODO implement me
	panic("implement me")
}

func (w *WorkoutRepo) Delete(id uint64) error {
	//TODO implement me
	panic("implement me")
}
