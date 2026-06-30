package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"1000words-game/internal/api"
	"1000words-game/internal/config"
	"1000words-game/internal/database"
	"1000words-game/internal/logger"
	"1000words-game/internal/repositories"
	"1000words-game/internal/services"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New()

	db, err := database.Connect(cfg.DBPath, appLogger)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)
	wordRepository := repositories.NewWordRepository(db)
	levelRepository := repositories.NewLevelRepository(db)
	languageRepository := repositories.NewLanguageRepository(db)
	attemptRepository := repositories.NewAttemptRepository(db)
	progressRepository := repositories.NewProgressRepository(db)
	authTokenRepository := repositories.NewAuthTokenRepository(db)
	authorizationRepository := repositories.NewAuthorizationRepository(db)

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

	authService := services.NewAuthService(
		userService,
		authTokenRepository,
		authorizationRepository,
		cfg.APITokenDays,
	)

	authorizationService := services.NewAuthorizationService(
		authorizationRepository,
	)

	router := api.NewRouter(api.RouterDependencies{
		DB: db,

		Config: cfg,

		AuthService:          authService,
		AuthorizationService: authorizationService,
		UserService:          userService,
		LanguageService:      languageService,
		ProgressService:      progressService,
		GameService:          gameService,
	})

	server := api.NewServer(cfg, router)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.Start()
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}

	case signalValue := <-shutdown:
		fmt.Printf("\nShutdown signal received: %s\n", signalValue.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		fmt.Println("API server stopped.")
	}
}
