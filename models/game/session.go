package game

type Session struct {
	AttemptID        int
	ProgressID       int
	UserID           int
	TargetLanguageID int
	NativeLanguageID int
	LevelID          int
	Questions        []Question
	Answers          []Answer
	MaxWrongAnswers  int
}
