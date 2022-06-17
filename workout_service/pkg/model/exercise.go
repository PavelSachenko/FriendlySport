package model

var ExerciseTable = "exercises"

type Exercise struct {
	ID          uint64 `json:"id" sql:"id" db:"id"`
	WorkoutId   uint64 `json:"workout_id" db:"workout_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	IsDone      bool   `json:"is_done" db:"is_done"`
}

type ExerciseIntoWorkout struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title" `
	Description string `json:"description" `
	IsDone      bool   `json:"is_done" `
}

type ExerciseUpdate struct {
	Id          uint64  `json:"id"`
	UserId      uint64  `json:"user_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsDone      *bool   `json:"is_done"`
}
