package cli

import (
	"fmt"
	"strconv"
	"strings"

	dbModels "1000words-game/models/db"
)

func (a *App) readText(label string) string {
	fmt.Print(label)

	value, _ := a.reader.ReadString('\n')
	value = strings.TrimSpace(value)

	return value
}

func (a *App) waitForEnter() {
	fmt.Println("Pritisni ENTER za povratak na meni...")
	_, _ = a.reader.ReadString('\n')
}

func (a *App) chooseLevel() (int, error) {
	levels, err := a.gameService.GetAvailableLevels()
	if err != nil {
		return 0, err
	}

	fmt.Println()
	fmt.Println("Izaberi level:")

	for _, level := range levels {
		fmt.Printf("%d. %s\n", level.LevelNumber, level.Name)
	}

	for {
		value := a.readText("Level: ")

		levelNumber, err := strconv.Atoi(value)
		if err != nil {
			printError("Unesi broj levela.")
			continue
		}

		for _, level := range levels {
			if level.LevelNumber == levelNumber {
				return levelNumber, nil
			}
		}

		printError("Taj level ne postoji.")
	}
}

func (a *App) chooseLanguages() (string, string, error) {
	languages, err := a.gameService.GetAvailableLanguages()
	if err != nil {
		return "", "", err
	}

	fmt.Println()
	fmt.Println("Dostupni jezici:")

	for _, language := range languages {
		fmt.Printf("%s - %s\n", language.Code, language.Name)
	}

	var targetCode string
	var nativeCode string

	for {
		targetCode = normalizeLanguageCode(
			a.readText("Target jezik koji učiš, npr. en/es/sr: "),
		)

		if languageCodeExists(languages, targetCode) {
			break
		}

		printError("Taj target jezik ne postoji.")
	}

	for {
		nativeCode = normalizeLanguageCode(
			a.readText("Native jezik, npr. sr/en/es: "),
		)

		if !languageCodeExists(languages, nativeCode) {
			printError("Taj native jezik ne postoji.")
			continue
		}

		if nativeCode == targetCode {
			printError("Target i native ne mogu biti isti.")
			continue
		}

		break
	}

	return targetCode, nativeCode, nil
}

func normalizeAnswer(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	return value
}

func normalizeLanguageCode(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	return value
}

func languageCodeExists(languages []dbModels.Language, code string) bool {
	code = normalizeLanguageCode(code)

	for _, language := range languages {
		if language.Code == code {
			return true
		}
	}

	return false
}
