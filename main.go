package main

import (
	"log"

	"1000words-game/internal/cli"
	"1000words-game/internal/config"
	"1000words-game/internal/database"
	"1000words-game/internal/logger"
	"1000words-game/internal/repositories"
	"1000words-game/internal/services"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New()

	appLogger.Info("Application started", "app", cfg.AppName)

	db, err := database.Connect(cfg.DBPath, appLogger)
	if err != nil {
		appLogger.Error("Database connection failed", "error", err)
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)
	wordRepository := repositories.NewWordRepository(db)
	levelRepository := repositories.NewLevelRepository(db)
	languageRepository := repositories.NewLanguageRepository(db)
	attemptRepository := repositories.NewAttemptRepository(db)
	progressRepository := repositories.NewProgressRepository(db)

	userService := services.NewUserService(userRepository)
	languageService := services.NewLanguageService(languageRepository)

	progressService := services.NewProgressService(
		progressRepository,
		languageRepository,
		levelRepository,
	)

	gameService := services.NewGameService(
		wordRepository,
		levelRepository,
		attemptRepository,
		progressRepository,
		languageService,
		cfg.WordsPerLevel,
	)

	cliApp := cli.NewApp(
		gameService,
		userService,
		progressService,
	)

	_ = cliApp.Run()

	appLogger.Info("Application finished")
}
