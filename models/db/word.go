package db

import "time"

type Word struct {
	ID        int
	LevelID   int
	CreatedAt time.Time
}
