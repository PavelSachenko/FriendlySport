package model

import "time"

var WorkoutTable = "workouts"

type Workout struct {
	ID            uint64    `json:"id" sql:"id" db:"id"`
	UserId        uint64    `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	IsDone        bool      `json:"is_done" db:"is_done"`
	AppointedTime time.Time `json:"appointed_time" db:"appointed_time"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type WorkoutRecommendation struct {
	Title string `json:"title"`
}
