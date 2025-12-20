package service

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	gamemodel "github.com/rajkovrga/1000words-game/internal/models/game"
	"github.com/rajkovrga/1000words-game/internal/repositories"
)

type SessionService struct {
	userRepostiory         *repositories.UserRepository
	wordRepository         *repositories.WordRepository
	levelResultsRepository *repositories.LevelResultsRepository
	languageRepository     *repositories.LanguageRepository
}

func NewGameService(
	userRepo *repositories.UserRepository,
	wordRepo *repositories.WordRepository,
	levelResultRepo *repositories.LevelResultsRepository,
	languageRepo *repositories.LanguageRepository,
) *SessionService {
	return &SessionService{
		userRepostiory:         userRepo,
		wordRepository:         wordRepo,
		levelResultsRepository: levelResultRepo,
		languageRepository:     languageRepo,
	}
}

func (sessionService *SessionService) Play() {
	user, err := sessionService.userRepostiory.GetFirstUser()
	language, err := sessionService.languageRepository.GetLanguagePerLanguageCode(os.Getenv("TARGET_LANGUAGE"))

	if err != nil {
		fmt.Println("User not found")
		os.Exit(0)
	}

	levelResult, err := sessionService.levelResultsRepository.GetLevelResultByUserId(user.Id)

	var wordPairs []gamemodel.WordPair

	defaultLevelStr := os.Getenv("DEFAULT_LEVEL")
	level, err := strconv.Atoi(defaultLevelStr)

	if err != nil {
		fmt.Println("Level default value not exist in env")
		os.Exit(0)
	}

	levelWordsStr := os.Getenv("LEVEL_WORDS")
	levelWords, err := strconv.Atoi(levelWordsStr)

	if err != nil {
		fmt.Println("Level default value not exist in env")
		os.Exit(0)
	}

	if levelResult != nil {
		level = levelResult.Level
	}

	wordPairs, _ = sessionService.wordRepository.GetGameWords(user.NativeLanguageId, language.Id, level, levelWords)

	maxMistakesStr := os.Getenv("MAX_MISTAKES")
	maxMistakes, err := strconv.Atoi(maxMistakesStr)

	mistakes := 0
	scanner := bufio.NewScanner(os.Stdin)

	for _, wp := range wordPairs {
		fmt.Printf("Translate '%s': ", wp.NativeWord)

		input := strings.TrimSpace(scanner.Text())

		if input != wp.TargetWord {
			fmt.Println("Wrong! Try again next time.")
			mistakes++
		} else {
			fmt.Println("Correct!")
		}

		if mistakes >= maxMistakes {
			fmt.Println("Max mistakes reached. Try again later.")
			break
		}
	}
}
