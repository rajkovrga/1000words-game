package services

import (
	"errors"
	"time"

	authUtil "1000words-game/internal/auth"
	"1000words-game/internal/repositories"
	dbModels "1000words-game/models/db"
)

type AuthService struct {
	userService             *UserService
	authTokenRepository     *repositories.AuthTokenRepository
	authorizationRepository *repositories.AuthorizationRepository
	tokenDays               int
}

type AuthResult struct {
	User      *dbModels.User
	Token     string
	ExpiresAt time.Time
}

func NewAuthService(
	userService *UserService,
	authTokenRepository *repositories.AuthTokenRepository,
	authorizationRepository *repositories.AuthorizationRepository,
	tokenDays int,
) *AuthService {
	if tokenDays <= 0 {
		tokenDays = 30
	}

	return &AuthService{
		userService:             userService,
		authTokenRepository:     authTokenRepository,
		authorizationRepository: authorizationRepository,
		tokenDays:               tokenDays,
	}
}

func (s *AuthService) Register(email string, password string) (*AuthResult, error) {
	user, err := s.userService.Register(email, password)
	if err != nil {
		return nil, err
	}

	if err := s.authorizationRepository.AssignRoleToUserByCode(user.ID, "user"); err != nil {
		return nil, err
	}

	return s.createTokenForUser(user)
}

func (s *AuthService) Login(email string, password string) (*AuthResult, error) {
	user, err := s.userService.Login(email, password)
	if err != nil {
		return nil, err
	}

	return s.createTokenForUser(user)
}

func (s *AuthService) Logout(rawToken string) error {
	if rawToken == "" {
		return errors.New("token is required")
	}

	tokenHash := authUtil.HashToken(rawToken)

	return s.authTokenRepository.DeleteByHash(tokenHash)
}

func (s *AuthService) Authenticate(rawToken string) (*dbModels.User, error) {
	if rawToken == "" {
		return nil, errors.New("token is required")
	}

	tokenHash := authUtil.HashToken(rawToken)

	token, err := s.authTokenRepository.FindValidByHash(tokenHash, time.Now())
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	user, err := s.userService.GetUserByID(token.UserID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	return user, nil
}

func (s *AuthService) createTokenForUser(user *dbModels.User) (*AuthResult, error) {
	rawToken, err := authUtil.GenerateToken()
	if err != nil {
		return nil, err
	}

	tokenHash := authUtil.HashToken(rawToken)
	expiresAt := time.Now().Add(time.Duration(s.tokenDays) * 24 * time.Hour)

	_, err = s.authTokenRepository.Create(
		user.ID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:      user,
		Token:     rawToken,
		ExpiresAt: expiresAt,
	}, nil
}
