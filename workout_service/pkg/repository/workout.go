package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "github.com/lib/pq"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/model"
	"log"
	"strings"
	"time"
)

type Workout interface {
	All(userId uint64, filterOption model.WorkoutsFiltering) (error, []*model.Workout)
	Create(workout *model.Workout) (error, *model.Workout)
	Update(workout model.WorkoutUpdate) (error, *model.Workout)
	Delete(id, userId uint64) error
	GetRecommendation(typingTitle string) (error, []*model.WorkoutRecommendation)
}

type WorkoutRepo struct {
	*db.DB
	elastic *elasticsearch.Client
}

func InitWorkoutRepo(db *db.DB, elClient *elasticsearch.Client) *WorkoutRepo {
	return &WorkoutRepo{
		DB:      db,
		elastic: elClient,
	}
}

type elasticSearchResults struct {
	Total int                            `json:"total"`
	Hits  []*model.WorkoutRecommendation `json:"hits"`
}

func (w *WorkoutRepo) GetRecommendation(typingTitle string) (error, []*model.WorkoutRecommendation) {
	err, res := w.searchByElastic(typingTitle)
	if err != nil {
		return err, nil
	}

	var results elasticSearchResults
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err, nil
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"]), nil
	}

	type envelopeResponse struct {
		Took int
		Hits struct {
			Total struct {
				Value int
			}
			Hits []struct {
				ID     string          `json:"_id"`
				Source json.RawMessage `json:"_source"`
				Title  string          `json:"title"`
			}
		}
	}
	var r envelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err, nil
	}
	results.Total = r.Hits.Total.Value
	for _, hit := range r.Hits.Hits {
		var wr model.WorkoutRecommendation
		if err := json.Unmarshal(hit.Source, &wr); err != nil {
			log.Fatalln(err)
		}
		results.Hits = append(results.Hits, &wr)
	}

	return nil, results.Hits
}

func (w *WorkoutRepo) searchByElastic(typingTitle string) (error, *esapi.Response) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query":     typingTitle,
					"fuzziness": 1,
				},
			},
		},
		"collapse": map[string]string{
			"field": "title.keyword",
		},
		"size": "5",
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err, nil
	}

	res, err := w.elastic.Search(
		w.elastic.Search.WithContext(context.Background()),
		w.elastic.Search.WithIndex("friendly_sport_workout_recommendation"),
		w.elastic.Search.WithBody(&buf),
		w.elastic.Search.WithPretty(),
	)
	if err != nil {
		return err, nil
	}

	return nil, res
}

func (w *WorkoutRepo) All(userId uint64, filterOption model.WorkoutsFiltering) (error, []*model.Workout) {
	sql := w.DB.QueryBuilder.Select("w.*").From(fmt.Sprintf("%s AS w", model.WorkoutTable)).Where("w.user_id = @", userId)
	if filterOption.Title != nil {
		sql = sql.AndWhere("w.title LIKE @", fmt.Sprint("%", filterOption.Title.(string), "%"))
	}
	if filterOption.IsDone != nil {
		sql = sql.AndWhere("w.is_done = @", filterOption.IsDone.(bool))
	}
	sql = sql.GroupBy("w.id")
	if filterOption.Sort != nil {
		sql = sql.OrderBy(strings.ReplaceAll(filterOption.Sort.(string), ":", " "))
	} else {
		sql = sql.OrderBy("created_at DESC")

	}
	sql = sql.Limit(filterOption.Limit)
	sql = sql.Offset(filterOption.Offset)
	var workouts []*model.Workout

	query, args := sql.ToSql()
	rows, err := w.DB.Query(query, args...)
	if err != nil {
		return err, nil
	}
	for rows.Next() {
		workout := model.Workout{}
		err := rows.Scan(&workout.ID, &workout.UserId, &workout.Title, &workout.Description, &workout.IsDone, &workout.AppointedTime, &workout.CreatedAt, &workout.UpdatedAt)
		if err != nil {
			return err, nil
		}
		sqlExercises, exercisesArgs := w.DB.QueryBuilder.Select("e.id, e.title, e.description, e.is_done").From(model.ExerciseTable+" e").
			Join(fmt.Sprintf("%s w ON w.id = e.workout_id", model.WorkoutTable)).
			Where("w.id = @", workout.ID).
			AndWhere("w.user_id = @", userId).ToSql()
		exercisesRows, err := w.DB.Query(sqlExercises, exercisesArgs...)
		if err != nil {
			return err, nil
		}
		var exercises []*model.ExerciseIntoWorkout
		for exercisesRows.Next() {
			exercise := model.ExerciseIntoWorkout{}
			err := exercisesRows.Scan(&exercise.ID, &exercise.Title, &exercise.Description, &exercise.IsDone)
			if err != nil {
				return err, nil
			}
			exercises = append(exercises, &exercise)
		}
		workout.Exercises = exercises
		workouts = append(workouts, &workout)
	}

	return nil, workouts
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
	if len(workout.Title) >= 3 {
		w.addToWorkoutTitleRecommendation(workout.Title)
	}

	return nil, workout
}

func (w *WorkoutRepo) Update(workoutUpdate model.WorkoutUpdate) (error, *model.Workout) {
	sql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update(model.WorkoutTable).Where("id = ?", workoutUpdate.Id).Where("user_id = ?", workoutUpdate.UserId)
	if workoutUpdate.IsDone != nil {
		sql = sql.Set("is_done", workoutUpdate.IsDone)
	}
	if workoutUpdate.Title != nil {
		sql = sql.Set("title", workoutUpdate.Title)
		if len(*workoutUpdate.Title) >= 3 {
			w.addToWorkoutTitleRecommendation(*workoutUpdate.Title)
		}
	}
	if workoutUpdate.Description != nil {
		sql = sql.Set("description", workoutUpdate.Description)
	}
	if workoutUpdate.AppointedTime != nil {
		utcTime := time.Unix(int64(*workoutUpdate.AppointedTime), 0).UTC()
		sql = sql.Set("appointed_time", utcTime)
	}
	sql = sql.Set("updated_at", workoutUpdate.UpdatedAt)
	query, args, _ := sql.ToSql()

	rows, err := w.DB.Queryx(query+" RETURNING *", args...)
	if err != nil {
		return err, nil
	}
	var workout model.Workout
	if rows.Next() {
		err = rows.Scan(&workout.ID, &workout.UserId, &workout.Title, &workout.Description, &workout.IsDone, &workout.AppointedTime, &workout.CreatedAt, &workout.UpdatedAt)
		if err != nil {
			return err, nil
		}
	}

	return nil, &workout
}

func (w *WorkoutRepo) Delete(id, userId uint64) error {
	res, err := w.DB.Exec("DELETE FROM "+model.WorkoutTable+" WHERE id=$1 AND user_id=$2", id, userId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count <= 0 {
		return errors.New("not found")
	}

	return nil
}

func (w *WorkoutRepo) addToWorkoutTitleRecommendation(title string) {
	//TODO add logger
	res, err := w.elastic.Index("friendly_sport_workout_recommendation", bytes.NewReader([]byte(fmt.Sprintf("{\"title\": \"%s\"}", title))))
	log.Println(err)
	log.Println(res)
}
