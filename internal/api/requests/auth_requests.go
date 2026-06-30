package requests

import (
	"errors"
	"strings"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(r.Password) == "" {
		return errors.New("password is required")
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(r.Password) == "" {
		return errors.New("password is required")
	}

	return nil
}
