package cli

import (
	"bufio"
	"os"

	"1000words-game/internal/services"
	dbModels "1000words-game/models/db"
)

type App struct {
	gameService     *services.GameService
	userService     *services.UserService
	progressService *services.ProgressService
	reader          *bufio.Reader
	currentUser     *dbModels.User
}

func NewApp(
	gameService *services.GameService,
	userService *services.UserService,
	progressService *services.ProgressService,
) *App {
	return &App{
		gameService:     gameService,
		userService:     userService,
		progressService: progressService,
		reader:          bufio.NewReader(os.Stdin),
		currentUser:     nil,
	}
}
