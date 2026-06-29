package services

import (
	"errors"
	"fmt"

	"1000words-game/internal/repositories"
	gameModels "1000words-game/models/game"
)

type ProgressService struct {
	progressRepository *repositories.ProgressRepository
	languageRepository *repositories.LanguageRepository
	levelRepository    *repositories.LevelRepository
}

func NewProgressService(
	progressRepository *repositories.ProgressRepository,
	languageRepository *repositories.LanguageRepository,
	levelRepository *repositories.LevelRepository,
) *ProgressService {
	return &ProgressService{
		progressRepository: progressRepository,
		languageRepository: languageRepository,
		levelRepository:    levelRepository,
	}
}

func (s *ProgressService) GetUserProgressOptions(userID int) ([]gameModels.ProgressOption, error) {
	if userID <= 0 {
		return nil, errors.New("user id is required")
	}

	return s.progressRepository.GetOptionsByUserID(userID)
}

func (s *ProgressService) CreateProgress(
	userID int,
	targetCode string,
	nativeCode string,
) (*gameModels.ProgressOption, error) {
	if userID <= 0 {
		return nil, errors.New("user id is required")
	}

	if targetCode == "" {
		return nil, errors.New("target language is required")
	}

	if nativeCode == "" {
		return nil, errors.New("native language is required")
	}

	if targetCode == nativeCode {
		return nil, errors.New("target and native language cannot be the same")
	}

	targetLanguage, err := s.languageRepository.FindByCode(targetCode)
	if err != nil {
		return nil, fmt.Errorf("target language not found: %w", err)
	}

	nativeLanguage, err := s.languageRepository.FindByCode(nativeCode)
	if err != nil {
		return nil, fmt.Errorf("native language not found: %w", err)
	}

	level, err := s.levelRepository.FindByLevelNumber(1)
	if err != nil {
		return nil, fmt.Errorf("level 1 not found: %w", err)
	}

	progress, err := s.progressRepository.Create(
		userID,
		targetLanguage.ID,
		nativeLanguage.ID,
		level.ID,
	)
	if err != nil {
		return nil, err
	}

	return s.progressRepository.FindOptionByID(progress.ID)
}
