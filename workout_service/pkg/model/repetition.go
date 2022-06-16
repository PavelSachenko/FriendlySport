package model

var RepetitionTable = "repetitions"

type Repetition struct {
	ID         uint64  `json:"id" sql:"id" db:"id"`
	ExerciseId uint64  `json:"exercise_id" db:"exercise_id"`
	Weight     float64 `json:"weight" db:"weight"`
	Count      uint64  `json:"count" db:"count"`
	IsDone     bool    `json:"is_done" db:"is_done"`
}

type RepetitionIntoExercise struct {
	ID     uint64  `json:"id" sql:"id" db:"id"`
	Weight float64 `json:"weight" db:"weight"`
	Count  uint64  `json:"count" db:"count"`
	IsDone bool    `json:"is_done" db:"is_done"`
}
