package db

import "time"

type Level struct {
	ID              int
	LevelNumber     int
	Name            string
	WordsRequired   int
	MaxWrongAnswers int
	CreatedAt       time.Time
}
