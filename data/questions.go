package data

import (
	"1000words-game/models"
)

func GetLevelWords(level int) []models.Question {
	return []models.Question{
		{Word: "cat", Answer: "mačka"},
		{Word: "dog", Answer: "pas"},
		{Word: "house", Answer: "kuća"},
		{Word: "car", Answer: "auto"},
	}
}
