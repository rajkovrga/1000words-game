package data

import (
	models "1000words-game/models/game"
)

func GetLevelWords(level int) []models.Question {
	return []models.Question{
		{Word: "cat", Answer: "mačka"},
		{Word: "dog", Answer: "pas"},
		{Word: "house", Answer: "kuća"},
		{Word: "car", Answer: "auto"},
	}
}
