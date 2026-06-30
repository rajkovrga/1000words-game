package requests

import "errors"

type StartGameRequest struct {
	ProgressID int `json:"progress_id"`
}

func (r StartGameRequest) Validate() error {
	if r.ProgressID <= 0 {
		return errors.New("progress_id is required")
	}

	return nil
}

type FinishGameRequest struct {
	AttemptID      int `json:"attempt_id"`
	ProgressID     int `json:"progress_id"`
	CorrectCount   int `json:"correct_count"`
	WrongCount     int `json:"wrong_count"`
	TotalQuestions int `json:"total_questions"`
}

func (r FinishGameRequest) Validate() error {
	if r.AttemptID <= 0 {
		return errors.New("attempt_id is required")
	}

	if r.ProgressID <= 0 {
		return errors.New("progress_id is required")
	}

	if r.CorrectCount < 0 {
		return errors.New("correct_count cannot be negative")
	}

	if r.WrongCount < 0 {
		return errors.New("wrong_count cannot be negative")
	}

	if r.TotalQuestions <= 0 {
		return errors.New("total_questions is required")
	}

	if r.CorrectCount+r.WrongCount != r.TotalQuestions {
		return errors.New("correct_count + wrong_count must be equal to total_questions")
	}

	return nil
}
