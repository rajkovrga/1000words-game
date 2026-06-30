package db

import "time"

type Language struct {
	ID        int
	Name      string
	Code      string
	CreatedAt time.Time
}
