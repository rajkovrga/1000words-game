package game

import dbModels "1000words-game/models/db"

type FinishGameResult struct {
	AttemptID       int
	ProgressID      int
	CorrectCount    int
	WrongCount      int
	TotalQuestions  int
	MaxWrongAnswers int
	Passed          bool
	NextLevel       *dbModels.Level
}
