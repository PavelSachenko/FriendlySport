package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "github.com/lib/pq"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/errors"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/model"
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
	logger  logger.Logging
}

func InitWorkoutRepo(logger logger.Logging, db *db.DB, elClient *elasticsearch.Client) *WorkoutRepo {
	return &WorkoutRepo{
		DB:      db,
		elastic: elClient,
		logger:  logger,
	}
}

type elasticSearchResults struct {
	Total int                            `json:"total"`
	Hits  []*model.WorkoutRecommendation `json:"hits"`
}

func (w *WorkoutRepo) GetRecommendation(typingTitle string) (error, []*model.WorkoutRecommendation) {
	err, res := w.searchByElastic(typingTitle)
	if err != nil {
		w.logger.Error(fmt.Sprintf("search recomendation from elastic error: ERROR %s", err.Error()))
		return errors.UnprocessableEntity, nil
	}

	var results elasticSearchResults
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			w.logger.Error(fmt.Sprintf("json decoder error: ERROR %s", err.Error()))
			return errors.UnprocessableEntity, nil
		}
		w.logger.Error(fmt.Sprintf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"]))

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
		return errors.UnprocessableEntity, nil
	}
	results.Total = r.Hits.Total.Value
	for _, hit := range r.Hits.Hits {
		var wr model.WorkoutRecommendation
		if err := json.Unmarshal(hit.Source, &wr); err != nil {
			w.logger.Error(fmt.Sprintf("json unmarshal error: ERROR %s", err.Error()))
			return errors.UnprocessableEntity, nil
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
		w.logger.Error(fmt.Sprintf("json encoder error: ERROR %s", err.Error()))
		return errors.UnprocessableEntity, nil
	}

	res, err := w.elastic.Search(
		w.elastic.Search.WithContext(context.Background()),
		w.elastic.Search.WithIndex("friendly_sport_workout_recommendation"),
		w.elastic.Search.WithBody(&buf),
		w.elastic.Search.WithPretty(),
	)
	if err != nil {
		w.logger.Error(fmt.Sprintf("elastic search error: ERROR %s", err.Error()))
		return errors.UnprocessableEntity, nil
	}

	return nil, res
}

func (w *WorkoutRepo) All(userId uint64, filterOption model.WorkoutsFiltering) (error, []*model.Workout) {

	query, args := w.createAllSqlRequest(userId, filterOption)
	rows, err := w.DB.Query(query, args...)
	if err != nil {
		w.logger.Error(fmt.Sprintf("query error: ERROR %s", err.Error()))
		return errors.UnprocessableEntity, nil
	}

	var workouts []*model.Workout
	for rows.Next() {
		workout := model.Workout{}
		var exercises []*model.ExerciseIntoWorkout
		var jsonExercises *string
		err := rows.Scan(&workout.ID, &workout.UserId, &workout.Title, &workout.Description, &workout.IsDone, &workout.AppointedTime, &workout.CreatedAt, &workout.UpdatedAt, &jsonExercises)
		if err != nil {
			w.logger.Error(fmt.Sprintf("row scan error: ERROR %s", err.Error()))
			return errors.UnprocessableEntity, nil
		}
		if jsonExercises != nil {
			err = json.Unmarshal([]byte(*jsonExercises), &exercises)
			if err != nil {
				w.logger.Error(fmt.Sprintf("unmarshal error: ERROR %s", err.Error()))
				return errors.UnprocessableEntity, nil
			}
			workout.Exercises = exercises
		}
		workouts = append(workouts, &workout)
	}

	return nil, workouts
}

func (w *WorkoutRepo) createAllSqlRequest(userId uint64, filterOption model.WorkoutsFiltering) (query string, arguments []interface{}) {
	sql := w.DB.QueryBuilder.Select("w.*").From(fmt.Sprintf("%s AS w", model.WorkoutView)).Where("w.user_id = @", userId)
	if filterOption.Title != nil {
		sql = sql.AndWhere("w.title LIKE @", fmt.Sprint("%", filterOption.Title.(string), "%"))
	}
	if filterOption.IsDone != nil {
		sql = sql.AndWhere("w.is_done = @", filterOption.IsDone.(bool))
	}
	if filterOption.Sort != nil {
		sql = sql.OrderBy(strings.ReplaceAll(filterOption.Sort.(string), ":", " "))
	} else {
		sql = sql.OrderBy("created_at DESC")

	}
	sql = sql.Limit(filterOption.Limit)
	sql = sql.Offset(filterOption.Offset)

	query, args := sql.ToSql()

	return query, args
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
			w.logger.Error(fmt.Sprintf("struct scan error: ERROR %s", err.Error()))
			return errors.UnprocessableEntity, nil
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
		return errors.UnprocessableEntity, nil
	}
	var workout model.Workout
	if rows.Next() {
		err = rows.Scan(&workout.ID, &workout.UserId, &workout.Title, &workout.Description, &workout.IsDone, &workout.AppointedTime, &workout.CreatedAt, &workout.UpdatedAt)
		if err != nil {
			w.logger.Error(fmt.Sprintf("row scan error: ERROR %s", err.Error()))
			return errors.UnprocessableEntity, nil
		}
	}

	return nil, &workout
}

func (w *WorkoutRepo) Delete(id, userId uint64) error {
	res, err := w.DB.Exec("DELETE FROM "+model.WorkoutTable+" WHERE id=$1 AND user_id=$2", id, userId)
	if err != nil {
		w.logger.Error(fmt.Sprintf("sql exec: ERROR %s", err.Error()))
		return errors.UnprocessableEntity
	}
	count, err := res.RowsAffected()
	if err != nil || count <= 0 {
		return errors.NotFound
	}

	return nil
}

func (w *WorkoutRepo) addToWorkoutTitleRecommendation(title string) {
	_, err := w.elastic.Index("friendly_sport_workout_recommendation", bytes.NewReader([]byte(fmt.Sprintf("{\"title\": \"%s\"}", title))))
	w.logger.Error(fmt.Sprintf("elastic add doc to indexerror: ERROR %s", err.Error()))

}
