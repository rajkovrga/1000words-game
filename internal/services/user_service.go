package services

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"1000words-game/internal/repositories"
	dbModels "1000words-game/models/db"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Register(email string, password string) (*dbModels.User, error) {
	email = normalizeEmail(email)

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	exists, err := s.userRepository.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.Create(email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email string, password string) (*dbModels.User, error) {
	email = normalizeEmail(email)

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *UserService) GetUserByID(userID int) (*dbModels.User, error) {
	if userID <= 0 {
		return nil, errors.New("user id must be greater than zero")
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (s *UserService) ChangePassword(userID int, newPassword string) error {
	if userID <= 0 {
		return errors.New("user id must be greater than zero")
	}

	if err := validatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepository.UpdatePassword(userID, hashedPassword)
}

func (s *UserService) DeleteUser(userID int) error {
	if userID <= 0 {
		return errors.New("user id must be greater than zero")
	}

	return s.userRepository.SoftDelete(userID)
}

func normalizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	return email
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if len(email) > 254 {
		return errors.New("email is too long")
	}

	if strings.ContainsAny(email, " \t\r\n") {
		return errors.New("email must not contain spaces")
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email is not valid")
	}

	if parsedEmail.Address != email {
		return errors.New("email is not valid")
	}

	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must have at least 8 characters")
	}

	if len(password) > 72 {
		return errors.New("password is too long")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func checkPasswordHash(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	return err == nil
}
