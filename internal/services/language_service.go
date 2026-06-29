package services

import (
	"errors"
	"fmt"

	"1000words-game/internal/repositories"
	dbModels "1000words-game/models/db"
)

type LanguageService struct {
	languageRepository *repositories.LanguageRepository
}

func NewLanguageService(languageRepository *repositories.LanguageRepository) *LanguageService {
	return &LanguageService{
		languageRepository: languageRepository,
	}
}

func (s *LanguageService) GetAllLanguages() ([]dbModels.Language, error) {
	return s.languageRepository.GetAll()
}

func (s *LanguageService) GetLanguageByID(id int) (*dbModels.Language, error) {
	if id <= 0 {
		return nil, errors.New("language id must be greater than zero")
	}

	language, err := s.languageRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("language not found by id %d: %w", id, err)
	}

	return language, nil
}

func (s *LanguageService) GetLanguageByCode(code string) (*dbModels.Language, error) {
	if code == "" {
		return nil, errors.New("language code is required")
	}

	language, err := s.languageRepository.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("language not found by code %s: %w", code, err)
	}

	return language, nil
}

func (s *LanguageService) ValidateLanguagePair(targetCode string, nativeCode string) error {
	if targetCode == "" {
		return errors.New("target language code is required")
	}

	if nativeCode == "" {
		return errors.New("native language code is required")
	}

	if targetCode == nativeCode {
		return errors.New("target language and native language cannot be the same")
	}

	targetExists, err := s.languageRepository.ExistsByCode(targetCode)
	if err != nil {
		return err
	}

	if !targetExists {
		return fmt.Errorf("target language does not exist: %s", targetCode)
	}

	nativeExists, err := s.languageRepository.ExistsByCode(nativeCode)
	if err != nil {
		return err
	}

	if !nativeExists {
		return fmt.Errorf("native language does not exist: %s", nativeCode)
	}

	return nil
}

func (s *LanguageService) GetLanguagePair(targetCode string, nativeCode string) (*dbModels.Language, *dbModels.Language, error) {
	err := s.ValidateLanguagePair(targetCode, nativeCode)
	if err != nil {
		return nil, nil, err
	}

	targetLanguage, err := s.languageRepository.FindByCode(targetCode)
	if err != nil {
		return nil, nil, err
	}

	nativeLanguage, err := s.languageRepository.FindByCode(nativeCode)
	if err != nil {
		return nil, nil, err
	}

	return targetLanguage, nativeLanguage, nil
}
