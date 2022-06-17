package repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/model"
	"log"
	"strings"
)

type Exercise interface {
	Create(exercise model.Exercise) (error error, createdExercise *model.Exercise)
	Delete(ctx context.Context) error
	Update(ctx context.Context, update model.ExerciseUpdate) (error, *model.Exercise)
	All(ctx context.Context) (error, []*model.ExerciseIntoWorkout)
}

type ExerciseRepo struct {
	*db.DB
}

func InitExerciseRepo(db *db.DB) ExerciseRepo {
	return ExerciseRepo{
		db,
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
			return err, nil
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
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count <= 0 {
		return errors.New("not found")
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
		return err, nil
	}
	var updatedExercise model.Exercise
	if rows.Next() {
		err := rows.Scan(&updatedExercise.ID, &updatedExercise.WorkoutId, &updatedExercise.Title, &updatedExercise.Description, &updatedExercise.IsDone)
		if err != nil {
			t.Rollback()
			return err, nil
		}
	}
	rows.NextResultSet()
	_, err = t.Exec(fmt.Sprintf("UPDATE %s SET updated_at = now() WHERE user_id = $1 AND id = $2", model.WorkoutTable), userId, workoutId)
	if err != nil {
		t.Rollback()
		return err, nil
	}
	err = t.Commit()
	if err != nil {
		return err, nil
	}

	return nil, &updatedExercise
}

func (e ExerciseRepo) All(ctx context.Context) (err error, exercises []*model.ExerciseIntoWorkout) {
	_, userId, workoutId := e.getIds(ctx)
	sql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("e.id, e.title, e.description, e.is_done").From(model.ExerciseTable+" as e").
		Join(model.WorkoutTable+" ON "+fmt.Sprintf("%s.id = e.workout_id", model.WorkoutTable)).
		Where(model.WorkoutTable+".user_id = ?", userId).
		Where("e.workout_id = ?", workoutId).
		OrderBy("e.id DESC")
	query, args, _ := sql.ToSql()
	fmt.Println(query)
	fmt.Println(args)
	rows, err := e.Query(query, args...)
	if err != nil {
		return err, nil
	}
	for rows.Next() {
		exercise := model.ExerciseIntoWorkout{}
		err := rows.Scan(&exercise.ID, &exercise.Title, &exercise.Description, &exercise.IsDone)
		if err != nil {
			return err, nil
		}
		exercises = append(exercises, &exercise)
	}
	return nil, exercises
}

func (e ExerciseRepo) getIds(ctx context.Context) (id, userId, workoutId uint64) {

	defer func() {
		//TODO add logger
		if r := recover(); r != nil {
			log.Fatalln("Converting error")
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
