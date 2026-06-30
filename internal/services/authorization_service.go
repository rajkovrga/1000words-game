package services

import (
	"errors"

	"1000words-game/internal/repositories"
)

type AuthorizationService struct {
	authorizationRepository *repositories.AuthorizationRepository
}

func NewAuthorizationService(
	authorizationRepository *repositories.AuthorizationRepository,
) *AuthorizationService {
	return &AuthorizationService{
		authorizationRepository: authorizationRepository,
	}
}

func (s *AuthorizationService) AssignUserRole(userID int) error {
	if userID <= 0 {
		return errors.New("user id is required")
	}

	return s.authorizationRepository.AssignRoleToUserByCode(userID, "user")
}

func (s *AuthorizationService) UserHasPermission(
	userID int,
	permissionCode string,
) (bool, error) {
	if userID <= 0 {
		return false, errors.New("user id is required")
	}

	if permissionCode == "" {
		return false, errors.New("permission code is required")
	}

	return s.authorizationRepository.UserHasPermission(userID, permissionCode)
}

func (s *AuthorizationService) GetUserRoleCodes(userID int) ([]string, error) {
	if userID <= 0 {
		return nil, errors.New("user id is required")
	}

	return s.authorizationRepository.GetUserRoleCodes(userID)
}

func (s *AuthorizationService) GetUserPermissionCodes(userID int) ([]string, error) {
	if userID <= 0 {
		return nil, errors.New("user id is required")
	}

	return s.authorizationRepository.GetUserPermissionCodes(userID)
}
