package db

import "time"

type UserLanguageProgress struct {
	ID                  int
	UserID              int
	TargetLanguageID    int
	NativeLanguageID    int
	CurrentLevelID      int
	TotalAttempts       int
	TotalPassed         int
	TotalFailed         int
	TotalCorrectAnswers int
	TotalWrongAnswers   int
	UpdatedAt           time.Time
	CreatedAt           time.Time
}
