package dbmodel

import "time"

type LevelResult struct {
	Id               int
	UserId           int
	TargetLanguageId int
	Level            int
	Success          bool
	LevelCompletedAt *time.Time
	TotalCompletedAt *time.Time
	CreatedAt        time.Time
}
