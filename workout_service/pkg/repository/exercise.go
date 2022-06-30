package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/errors"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/model"
	"github.com/pavel/workout_service/utils"
	"strings"
)

type Exercise interface {
	Create(userId uint64, exercise model.Exercise) (error error, createdExercise *model.Exercise)
	Delete(id, userId, workoutId uint64) error
	Update(update model.ExerciseUpdate) (error, *model.Exercise)
	GetRecommendation(typingTitle string) (error, []*model.ExerciseRecommendation)
}

type ExercisePostgresRepo struct {
	*db.DB
	elastic *elasticsearch.Client
	logger  *logger.Logger
}

func InitExercisePostgresRepo(logger *logger.Logger, db *db.DB, elClient *elasticsearch.Client) *ExercisePostgresRepo {
	return &ExercisePostgresRepo{
		DB:      db,
		logger:  logger,
		elastic: elClient,
	}
}

func (e *ExercisePostgresRepo) Create(userId uint64, exercise model.Exercise) (error, *model.Exercise) {

	res, err := e.Exec(fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1 AND id = $2", model.WorkoutTable), userId, exercise.WorkoutId)
	if err != nil {
		e.logger.Error(err.Error())
		return err, nil
	}
	isWorkoutExist, err := res.RowsAffected()
	if err != nil || isWorkoutExist <= 0 {
		e.logger.Error(err.Error())
		return errors.NotFound, nil
	}

	sql := fmt.Sprintf("INSERT INTO %s (workout_id, title, description) ", model.ExerciseTable)
	rows, err := e.Queryx(sql+"VALUES ($1, $2, $3) RETURNING *",
		exercise.WorkoutId,
		exercise.Title,
		exercise.Description,
	)
	if err != nil {
		e.logger.Error(err.Error())
		return err, nil
	}
	var createdExercise model.Exercise
	for rows.Next() {
		err = rows.StructScan(&createdExercise)
		if err != nil {
			e.logger.Errorf("struct scan: ERROR %s", err.Error())
			return errors.UnprocessableEntity, nil
		}
	}
	if len(createdExercise.Title) >= 3 {
		e.addToExerciseTitleRecommendation(createdExercise.Title)
	}

	return nil, &createdExercise
}

func (e ExercisePostgresRepo) Delete(id, userId, workoutId uint64) error {
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

func (e ExercisePostgresRepo) Update(update model.ExerciseUpdate) (error, *model.Exercise) {
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
	args = append(args, update.Id, update.UserId, update.WorkoutId)
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
	_, err = t.Exec(fmt.Sprintf("UPDATE %s SET updated_at = now() WHERE user_id = $1 AND id = $2", model.WorkoutTable), update.UserId, update.WorkoutId)
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

	if len(updatedExercise.Title) >= 3 {
		e.addToExerciseTitleRecommendation(updatedExercise.Title)
	}

	return nil, &updatedExercise
}

type elasticSearchExerciseRecommendationResults struct {
	Total int                             `json:"total"`
	Hits  []*model.ExerciseRecommendation `json:"hits"`
}

func (e *ExercisePostgresRepo) GetRecommendation(typingTitle string) (error, []*model.ExerciseRecommendation) {
	err, res := e.searchRecommendationTitleByElastic(typingTitle)
	if err != nil {
		e.logger.Errorf("search recomendation from elastic error: ERROR %s", err.Error())
		return errors.UnprocessableEntity, nil
	}

	err = utils.ValidateElasticResponse(res, e.logger)
	if err != nil {
		return err, nil
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
	var results elasticSearchExerciseRecommendationResults
	results.Total = r.Hits.Total.Value
	for _, hit := range r.Hits.Hits {
		var er model.ExerciseRecommendation
		if err := json.Unmarshal(hit.Source, &er); err != nil {
			e.logger.Errorf("json unmarshal error: ERROR %s", err.Error())
			return errors.UnprocessableEntity, nil
		}
		results.Hits = append(results.Hits, &er)
	}

	return nil, results.Hits
}

func (w *ExercisePostgresRepo) searchRecommendationTitleByElastic(typingTitle string) (error, *esapi.Response) {
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
		w.elastic.Search.WithIndex("friendly_sport_exercise_recommendation"),
		w.elastic.Search.WithBody(&buf),
		w.elastic.Search.WithPretty(),
	)
	if err != nil {
		w.logger.Errorf("elastic search error: ERROR %s", err.Error())
		return errors.UnprocessableEntity, nil
	}

	return nil, res
}

func (e ExercisePostgresRepo) addToExerciseTitleRecommendation(title string) {
	_, err := e.elastic.Index("friendly_sport_exercise_recommendation", bytes.NewReader([]byte(fmt.Sprintf("{\"title\": \"%s\"}", title))))
	if err != nil {
		e.logger.Errorf("elastic add doc to indexerror: ERROR %s", err.Error())
	}

}
