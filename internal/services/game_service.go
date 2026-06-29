package services

import (
	"database/sql"
	"errors"

	"1000words-game/internal/repositories"
	dbModels "1000words-game/models/db"
	gameModels "1000words-game/models/game"
)

type GameService struct {
	wordRepository     *repositories.WordRepository
	levelRepository    *repositories.LevelRepository
	attemptRepository  *repositories.AttemptRepository
	progressRepository *repositories.ProgressRepository
	languageService    *LanguageService
	wordsPerLevel      int
}

func NewGameService(
	wordRepository *repositories.WordRepository,
	levelRepository *repositories.LevelRepository,
	attemptRepository *repositories.AttemptRepository,
	progressRepository *repositories.ProgressRepository,
	languageService *LanguageService,
	wordsPerLevel int,
) *GameService {
	return &GameService{
		wordRepository:     wordRepository,
		levelRepository:    levelRepository,
		attemptRepository:  attemptRepository,
		progressRepository: progressRepository,
		languageService:    languageService,
		wordsPerLevel:      wordsPerLevel,
	}
}

func (s *GameService) GetAvailableLevels() ([]dbModels.Level, error) {
	return s.levelRepository.GetAll()
}

func (s *GameService) GetAvailableLanguages() ([]dbModels.Language, error) {
	return s.languageService.GetAllLanguages()
}

func (s *GameService) StartPractice(
	levelNumber int,
	targetLanguageCode string,
	nativeLanguageCode string,
) (*gameModels.Session, error) {
	targetLanguage, nativeLanguage, err := s.languageService.GetLanguagePair(
		targetLanguageCode,
		nativeLanguageCode,
	)
	if err != nil {
		return nil, err
	}

	level, err := s.levelRepository.FindByLevelNumber(levelNumber)
	if err != nil {
		return nil, err
	}

	questions, err := s.wordRepository.GetQuestionsByLevel(
		level.ID,
		targetLanguage.Code,
		nativeLanguage.Code,
		s.wordsPerLevel,
	)
	if err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, errors.New("no questions found for selected level and languages")
	}

	return &gameModels.Session{
		AttemptID:        0,
		ProgressID:       0,
		UserID:           0,
		TargetLanguageID: targetLanguage.ID,
		NativeLanguageID: nativeLanguage.ID,
		LevelID:          level.ID,
		Questions:        questions,
		Answers:          []gameModels.Answer{},
		MaxWrongAnswers:  level.MaxWrongAnswers,
	}, nil
}

func (s *GameService) StartLevel(
	userID int,
	progress gameModels.ProgressOption,
) (*gameModels.Session, error) {
	if userID <= 0 {
		return nil, errors.New("user id is required")
	}

	questions, err := s.wordRepository.GetQuestionsByLevel(
		progress.CurrentLevelID,
		progress.TargetCode,
		progress.NativeCode,
		s.wordsPerLevel,
	)
	if err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, errors.New("no questions found for selected user progress")
	}

	level, err := s.levelRepository.FindByID(progress.CurrentLevelID)
	if err != nil {
		return nil, err
	}

	attempt, err := s.attemptRepository.Create(
		userID,
		progress.TargetLanguageID,
		progress.NativeLanguageID,
		progress.CurrentLevelID,
		len(questions),
	)
	if err != nil {
		return nil, err
	}

	return &gameModels.Session{
		AttemptID:        attempt.ID,
		ProgressID:       progress.ProgressID,
		UserID:           userID,
		TargetLanguageID: progress.TargetLanguageID,
		NativeLanguageID: progress.NativeLanguageID,
		LevelID:          progress.CurrentLevelID,
		Questions:        questions,
		Answers:          []gameModels.Answer{},
		MaxWrongAnswers:  level.MaxWrongAnswers,
	}, nil
}

func (s *GameService) FinishLevel(
	attemptID int,
	progressID int,
	result gameModels.Result,
) (*dbModels.Level, error) {
	if attemptID <= 0 {
		return nil, errors.New("attempt id is required")
	}

	attempt, err := s.attemptRepository.Finish(attemptID, result)
	if err != nil {
		return nil, err
	}

	nextLevelID := 0
	var nextLevel *dbModels.Level

	if result.Passed {
		nextLevel, err = s.levelRepository.GetNextLevel(attempt.LevelID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if nextLevel != nil {
			nextLevelID = nextLevel.ID
		}
	}

	if progressID > 0 {
		err = s.progressRepository.UpdateAfterAttempt(
			progressID,
			result,
			nextLevelID,
		)
		if err != nil {
			return nil, err
		}
	}

	return nextLevel, nil
}
