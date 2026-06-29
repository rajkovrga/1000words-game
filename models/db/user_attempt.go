package db

import "time"

type UserAttempt struct {
	ID               int
	UserID           int
	TargetLanguageID int
	NativeLanguageID int
	LevelID          int
	CorrectCount     int
	WrongCount       int
	TotalQuestions   int
	Passed           bool
	StartedAt        time.Time
	FinishedAt       *time.Time
}
