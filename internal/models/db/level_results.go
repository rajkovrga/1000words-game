package dbmodel

import (
	"database/sql"
	"time"
)

type LevelResult struct {
	Id               int
	UserId           int
	TargetLanguageId int
	Level            int
	Success          bool
	LevelCompletedAt sql.NullTime
	TotalCompletedAt sql.NullTime
	CreatedAt        time.Time
}
