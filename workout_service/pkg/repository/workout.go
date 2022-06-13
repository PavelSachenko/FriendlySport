package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "github.com/lib/pq"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/model"
	"log"
	"time"
)

type Workout interface {
	One(id uint64) (error, *model.Workout)
	All()
	Create(workout *model.Workout) (error, *model.Workout)
	Update(workout *model.Workout) (error, *model.Workout)
	Delete(id uint64) error
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

	//TODO add logger
	data, err := json.Marshal(struct {
		Title string `json:"title"`
	}{Title: workout.Title})
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	w.elastic.Index("friendly_sport_workout_recommendation", bytes.NewReader(data))

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
