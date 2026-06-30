package requests

import (
	"errors"
	"strings"
)

type StartPracticeRequest struct {
	LevelNumber        int    `json:"level_number"`
	TargetLanguageCode string `json:"target_language_code"`
	NativeLanguageCode string `json:"native_language_code"`
}

func (r StartPracticeRequest) Validate() error {
	if r.LevelNumber <= 0 {
		return errors.New("level_number is required")
	}

	if strings.TrimSpace(r.TargetLanguageCode) == "" {
		return errors.New("target_language_code is required")
	}

	if strings.TrimSpace(r.NativeLanguageCode) == "" {
		return errors.New("native_language_code is required")
	}

	if strings.EqualFold(
		strings.TrimSpace(r.TargetLanguageCode),
		strings.TrimSpace(r.NativeLanguageCode),
	) {
		return errors.New("target_language_code and native_language_code cannot be the same")
	}

	return nil
}
