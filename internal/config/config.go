package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName       string
	DBPath        string
	WordsPerLevel int
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppName:       getEnv("APP_NAME", "1000words-game"),
		DBPath:        getEnv("DB_PATH", "words.db"),
		WordsPerLevel: getEnvAsInt("WORDS_PER_LEVEL", 100),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
