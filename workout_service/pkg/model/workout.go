package model

import (
	"encoding/json"
	"time"
)

var WorkoutTable = "workouts"
var WorkoutView = "workouts_objects"

type Workout struct {
	ID            uint64                 `json:"id" sql:"id" db:"id"`
	UserId        uint64                 `json:"user_id" db:"user_id"`
	Title         string                 `json:"title" db:"title"`
	Description   string                 `json:"description" db:"description"`
	IsDone        bool                   `json:"is_done" db:"is_done"`
	AppointedTime *time.Time             `json:"appointed_time" db:"appointed_time"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
	Exercises     []*ExerciseIntoWorkout `json:"exercises"`
}

func (u *Workout) MarshalJSON() ([]byte, error) {
	type workout Workout
	var unixAppointedTime int64
	if u.AppointedTime != nil {
		unixAppointedTime = u.AppointedTime.UTC().Unix()
	}
	return json.Marshal(&struct {
		*workout
		UpdatedAt     int64 `json:"updated_at"`
		CreatedAt     int64 `json:"created_at"`
		AppointedTime int64 `json:"appointed_time"`
	}{
		workout:       (*workout)(u),
		UpdatedAt:     u.UpdatedAt.UTC().Unix(),
		CreatedAt:     u.CreatedAt.UTC().Unix(),
		AppointedTime: unixAppointedTime,
	})

}

type WorkoutsFiltering struct {
	Title         interface{} `json:"title"`
	IsDone        interface{} `json:"is_done"`
	AppointedTime interface{} `json:"appointed_time"`
	Sort          interface{} `json:"sort"`
	Offset        uint64      `json:"offset"`
	Limit         uint64      `json:"limit"`
}

type WorkoutUpdate struct {
	Id            uint64    `json:"id"`
	UserId        uint64    `json:"user_id"`
	Title         *string   `json:"title"`
	Description   *string   `json:"description"`
	IsDone        *bool     `json:"is_done"`
	AppointedTime *uint64   `json:"appointed_time"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WorkoutRecommendation struct {
	Title string `json:"title"`
}
