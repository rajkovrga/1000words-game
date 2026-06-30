package routes

import (
	"database/sql"

	"1000words-game/internal/config"
	"1000words-game/internal/services"
)

type Dependencies struct {
	DB *sql.DB

	Config config.Config

	AuthService          *services.AuthService
	AuthorizationService *services.AuthorizationService
	UserService          *services.UserService
	LanguageService      *services.LanguageService
	ProgressService      *services.ProgressService
	GameService          *services.GameService
}
