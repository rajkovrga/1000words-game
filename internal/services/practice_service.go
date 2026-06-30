package services

import (
	"errors"

	"1000words-game/internal/repositories"
	gameModels "1000words-game/models/game"
)

type PracticeService struct {
	wordRepository  *repositories.WordRepository
	levelRepository *repositories.LevelRepository
	languageService *LanguageService
	wordsPerLevel   int
}

func NewPracticeService(
	wordRepository *repositories.WordRepository,
	levelRepository *repositories.LevelRepository,
	languageService *LanguageService,
	wordsPerLevel int,
) *PracticeService {
	return &PracticeService{
		wordRepository:  wordRepository,
		levelRepository: levelRepository,
		languageService: languageService,
		wordsPerLevel:   wordsPerLevel,
	}
}

func (s *PracticeService) StartPractice(
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
