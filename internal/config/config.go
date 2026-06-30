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

	APIPort                string
	APIReadTimeoutSeconds  int
	APIWriteTimeoutSeconds int
	APIIdleTimeoutSeconds  int
	APITokenDays           int
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppName:       getEnv("APP_NAME", "1000words-game"),
		DBPath:        getEnv("DB_PATH", "words.db"),
		WordsPerLevel: getEnvAsInt("WORDS_PER_LEVEL", 100),

		APIPort:                getEnv("API_PORT", "8080"),
		APIReadTimeoutSeconds:  getEnvAsInt("API_READ_TIMEOUT_SECONDS", 10),
		APIWriteTimeoutSeconds: getEnvAsInt("API_WRITE_TIMEOUT_SECONDS", 10),
		APIIdleTimeoutSeconds:  getEnvAsInt("API_IDLE_TIMEOUT_SECONDS", 60),
		APITokenDays:           getEnvAsInt("API_TOKEN_DAYS", 30),
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
