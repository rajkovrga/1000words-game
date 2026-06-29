package db

import "time"

type WordLanguage struct {
	ID             int
	WordID         int
	LanguageID     int
	Text           string
	NormalizedText string
	CreatedAt      time.Time
}
