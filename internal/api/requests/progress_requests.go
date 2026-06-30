package requests

import (
	"errors"
	"strings"
)

type CreateProgressRequest struct {
	TargetLanguageCode string `json:"target_language_code"`
	NativeLanguageCode string `json:"native_language_code"`
}

func (r CreateProgressRequest) Validate() error {
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
