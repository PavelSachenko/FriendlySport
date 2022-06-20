package repository

import (
	"context"
	"fmt"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/errors"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/model"
	"strings"
)

type Exercise interface {
	Create(exercise model.Exercise) (error error, createdExercise *model.Exercise)
	Delete(ctx context.Context) error
	Update(ctx context.Context, update model.ExerciseUpdate) (error, *model.Exercise)
}

type ExerciseRepo struct {
	*db.DB
	logger logger.Logger
}

func InitExerciseRepo(logger logger.Logger, db *db.DB) ExerciseRepo {
	return ExerciseRepo{
		DB:     db,
		logger: logger,
	}
}
func (e *ExerciseRepo) Create(exercise model.Exercise) (error error, createdExercise *model.Exercise) {
	sql := fmt.Sprintf("INSERT INTO %s (workout_id, title, description) ", model.ExerciseTable)
	rows, err := e.Queryx(sql+"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		exercise.WorkoutId,
		exercise.Title,
		exercise.Description,
	)
	if err != nil {
		return err, nil
	}
	for rows.Next() {
		err = rows.StructScan(&createdExercise)
		if err != nil {
			e.logger.Error(fmt.Sprintf("struct scan: ERROR %s", err.Error()))
			return errors.UnprocessableEntity, nil
		}
	}

	return nil, createdExercise
}

func (e ExerciseRepo) Delete(ctx context.Context) error {

	id, userId, workoutId := e.getIds(ctx)
	res, err := e.Exec(fmt.Sprintf("DELETE FROM %s USING %s WHERE %s.user_id = $1 AND %s.id = $2 AND %s.id = $3",
		model.ExerciseTable,
		model.WorkoutTable,
		model.WorkoutTable,
		model.WorkoutTable,
		model.ExerciseTable,
	), userId, workoutId, id)
	if err != nil {
		e.logger.Error(fmt.Sprintf("sql exec: ERROR %s", err.Error()))
		return errors.UnprocessableEntity
	}
	count, err := res.RowsAffected()
	if err != nil || count <= 0 {
		return errors.NotFound
	}

	return nil
}

func (e ExerciseRepo) Update(ctx context.Context, update model.ExerciseUpdate) (error, *model.Exercise) {
	id, userId, workoutId := e.getIds(ctx)
	args := make([]interface{}, 0)
	sets := make([]string, 0)
	argId := 1
	if update.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argId))
		args = append(args, update.Description)
		argId++
	}
	if update.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", argId))
		args = append(args, update.Title)
		argId++
	}
	if update.IsDone != nil {
		sets = append(sets, fmt.Sprintf("is_done = $%d", argId))
		args = append(args, update.IsDone)
		argId++
	}
	setQuery := strings.Join(sets, ",")
	sql := fmt.Sprintf("UPDATE %s as e SET %s FROM %s w WHERE w.id = e.workout_id AND e.id = $%d AND w.user_id = $%d AND e.workout_id = $%d RETURNING e.*",
		model.ExerciseTable, setQuery, model.WorkoutTable, argId, argId+1, argId+2)
	args = append(args, id, userId, workoutId)
	t, _ := e.Begin()
	rows, err := t.Query(sql, args...)
	if err != nil {
		t.Rollback()
		e.logger.Error(fmt.Sprintf("query prepare: ERROR %s", err.Error()))
		return errors.UnprocessableEntity, nil
	}
	var updatedExercise model.Exercise
	if rows.Next() {
		err := rows.Scan(&updatedExercise.ID, &updatedExercise.WorkoutId, &updatedExercise.Title, &updatedExercise.Description, &updatedExercise.IsDone)
		if err != nil {
			t.Rollback()
			e.logger.Error(fmt.Sprintf("rows scan: ERROR %s", err.Error()))
			return err, nil
		}
	}
	rows.NextResultSet()
	_, err = t.Exec(fmt.Sprintf("UPDATE %s SET updated_at = now() WHERE user_id = $1 AND id = $2", model.WorkoutTable), userId, workoutId)
	if err != nil {
		t.Rollback()
		e.logger.Error(fmt.Sprintf("sql exec: ERROR %s", err.Error()))
		return errors.UnprocessableEntity, nil
	}
	err = t.Commit()
	if err != nil {
		e.logger.Error(fmt.Sprintf("transaction commit: ERROR %s", err.Error()))
		return errors.UnprocessableEntity, nil
	}

	return nil, &updatedExercise
}

func (e ExerciseRepo) getIds(ctx context.Context) (id, userId, workoutId uint64) {

	defer func() {
		if r := recover(); r != nil {
			e.logger.Error(fmt.Sprintf("Converting error: ERROR %v", r))
		}
	}()

	if ctx.Value("id") != nil {
		id = ctx.Value("id").(uint64)
	}
	if ctx.Value("user_id") != nil {
		userId = ctx.Value("user_id").(uint64)
	}
	if ctx.Value("workout_id") != nil {
		workoutId = ctx.Value("workout_id").(uint64)
	}

	return id, userId, workoutId
}
