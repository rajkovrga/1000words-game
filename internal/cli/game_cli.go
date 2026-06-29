package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"1000words-game/internal/services"
	dbModels "1000words-game/models/db"
	gameModels "1000words-game/models/game"
	"1000words-game/utils"
)

type GameCLI struct {
	gameService     *services.GameService
	userService     *services.UserService
	progressService *services.ProgressService
	reader          *bufio.Reader
	currentUser     *dbModels.User
}

func NewGameCLI(
	gameService *services.GameService,
	userService *services.UserService,
	progressService *services.ProgressService,
) *GameCLI {
	return &GameCLI{
		gameService:     gameService,
		userService:     userService,
		progressService: progressService,
		reader:          bufio.NewReader(os.Stdin),
		currentUser:     nil,
	}
}

func (c *GameCLI) Run() error {
	for {
		fmt.Println("======================================")
		fmt.Println("1000words game")
		fmt.Println("======================================")
		fmt.Println("1. Praktika")
		fmt.Println("2. Igra")
		fmt.Println("0. Izlaz")
		fmt.Println("--------------------------------------")

		choice := c.readText("Izaberi opciju: ")

		switch choice {
		case "1":
			if err := c.runPractice(); err != nil {
				fmt.Printf("%sGreška: %s%s\n\n", utils.Red, err.Error(), utils.Reset)
				c.waitForEnter()
			}

		case "2":
			if err := c.runGame(); err != nil {
				fmt.Printf("%sGreška: %s%s\n\n", utils.Red, err.Error(), utils.Reset)
				c.waitForEnter()
			}

		case "0":
			fmt.Println("Izlaz iz aplikacije.")
			return nil

		default:
			fmt.Printf("%sNepoznata opcija.%s\n\n", utils.Red, utils.Reset)
		}
	}
}

func (c *GameCLI) runPractice() error {
	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("PRAKTIKA")
	fmt.Println("======================================")

	levelNumber, err := c.chooseLevel()
	if err != nil {
		return err
	}

	targetCode, nativeCode, err := c.chooseLanguages()
	if err != nil {
		return err
	}

	session, err := c.gameService.StartPractice(
		levelNumber,
		targetCode,
		nativeCode,
	)
	if err != nil {
		return err
	}

	if err := c.runQuestions(session, false, false); err != nil {
		return err
	}

	c.waitForEnter()
	return nil
}

func (c *GameCLI) runGame() error {
	if c.currentUser != nil {
		return c.startGameForUser(c.currentUser)
	}

	for {
		fmt.Println()
		fmt.Println("======================================")
		fmt.Println("IGRA")
		fmt.Println("======================================")
		fmt.Println("1. Prijava")
		fmt.Println("2. Registracija")
		fmt.Println("0. Nazad")
		fmt.Println("--------------------------------------")

		choice := c.readText("Izaberi opciju: ")

		switch choice {
		case "1":
			user, err := c.login()
			if err != nil {
				return err
			}

			c.currentUser = user

			return c.startGameForUser(c.currentUser)

		case "2":
			user, err := c.register()
			if err != nil {
				return err
			}

			c.currentUser = user

			return c.startGameForUser(c.currentUser)

		case "0":
			return nil

		default:
			fmt.Printf("%sNepoznata opcija.%s\n", utils.Red, utils.Reset)
		}
	}
}
func (c *GameCLI) startGameForUser(user *dbModels.User) error {
	fmt.Println()
	fmt.Printf("%sKorisnik: %s%s\n", utils.Green, user.Email, utils.Reset)

	progress, err := c.chooseOrCreateProgress(user.ID)
	if err != nil {
		return err
	}

	fmt.Printf(
		"%sIgra: %s -> %s | Level %d%s\n\n",
		utils.Blue,
		progress.TargetName,
		progress.NativeName,
		progress.LevelNumber,
		utils.Reset,
	)

	session, err := c.gameService.StartLevel(
		user.ID,
		*progress,
	)
	if err != nil {
		return err
	}

	if err := c.runQuestions(session, true, true); err != nil {
		return err
	}

	c.waitForEnter()
	return nil
}

func (c *GameCLI) login() (*dbModels.User, error) {
	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("PRIJAVA")
	fmt.Println("======================================")

	email := c.readText("Email: ")
	password := c.readText("Password: ")

	user, err := c.userService.Login(email, password)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%sUspešno ste se prijavili.%s\n", utils.Green, utils.Reset)
	return user, nil
}

func (c *GameCLI) register() (*dbModels.User, error) {
	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("REGISTRACIJA")
	fmt.Println("======================================")

	email := c.readText("Email: ")
	password := c.readText("Password: ")

	user, err := c.userService.Register(email, password)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%sUspešno ste se registrovali.%s\n", utils.Green, utils.Reset)
	return user, nil
}

func (c *GameCLI) chooseOrCreateProgress(userID int) (*gameModels.ProgressOption, error) {
	options, err := c.progressService.GetUserProgressOptions(userID)
	if err != nil {
		return nil, err
	}

	if len(options) == 0 {
		fmt.Println()
		fmt.Println("Nemaš podešen jezik za igru.")
		fmt.Println("Prvo izaberi target i native jezik.")

		targetCode, nativeCode, err := c.chooseLanguages()
		if err != nil {
			return nil, err
		}

		return c.progressService.CreateProgress(
			userID,
			targetCode,
			nativeCode,
		)
	}

	if len(options) == 1 {
		return &options[0], nil
	}

	fmt.Println()
	fmt.Println("Izaberi napredak/jezik za igru:")

	for index, option := range options {
		fmt.Printf(
			"%d. %s -> %s | Level %d\n",
			index+1,
			option.TargetName,
			option.NativeName,
			option.LevelNumber,
		)
	}

	for {
		value := c.readText("Opcija: ")

		selectedIndex, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("%sUnesi broj opcije.%s\n", utils.Red, utils.Reset)
			continue
		}

		if selectedIndex < 1 || selectedIndex > len(options) {
			fmt.Printf("%sTa opcija ne postoji.%s\n", utils.Red, utils.Reset)
			continue
		}

		return &options[selectedIndex-1], nil
	}
}

func (c *GameCLI) chooseLevel() (int, error) {
	levels, err := c.gameService.GetAvailableLevels()
	if err != nil {
		return 0, err
	}

	fmt.Println()
	fmt.Println("Izaberi level:")
	for _, level := range levels {
		fmt.Printf("%d. %s\n", level.LevelNumber, level.Name)
	}

	for {
		value := c.readText("Level: ")

		levelNumber, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("%sUnesi broj levela.%s\n", utils.Red, utils.Reset)
			continue
		}

		for _, level := range levels {
			if level.LevelNumber == levelNumber {
				return levelNumber, nil
			}
		}

		fmt.Printf("%sTaj level ne postoji.%s\n", utils.Red, utils.Reset)
	}
}

func (c *GameCLI) chooseLanguages() (string, string, error) {
	languages, err := c.gameService.GetAvailableLanguages()
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
			c.readText("Target jezik koji učiš, npr. en/es/sr: "),
		)

		if languageCodeExists(languages, targetCode) {
			break
		}

		fmt.Printf("%sTaj target jezik ne postoji.%s\n", utils.Red, utils.Reset)
	}

	for {
		nativeCode = normalizeLanguageCode(
			c.readText("Native jezik, npr. sr/en/es: "),
		)

		if !languageCodeExists(languages, nativeCode) {
			fmt.Printf("%sTaj native jezik ne postoji.%s\n", utils.Red, utils.Reset)
			continue
		}

		if nativeCode == targetCode {
			fmt.Printf("%sTarget i native ne mogu biti isti.%s\n", utils.Red, utils.Reset)
			continue
		}

		break
	}

	return targetCode, nativeCode, nil
}

func (c *GameCLI) runQuestions(
	session *gameModels.Session,
	saveResult bool,
	stopOnMaxWrong bool,
) error {
	correctCount := 0
	wrongCount := 0

	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("KVIZ POČINJE")
	fmt.Println("======================================")
	fmt.Printf("Ukupno pitanja: %d\n", len(session.Questions))

	if stopOnMaxWrong {
		fmt.Printf("Dozvoljene greške: %d\n", session.MaxWrongAnswers)
	} else {
		fmt.Println("Mod: praktika, nema ispadanja.")
	}

	fmt.Printf("%sZa izlaz iz kviza ukucaj #%s\n", utils.Yellow, utils.Reset)
	fmt.Println("--------------------------------------")
	fmt.Println()

	for index, question := range session.Questions {
		fmt.Printf("Pitanje %d/%d\n", index+1, len(session.Questions))
		fmt.Printf("Reč: %s\n", question.Word)

		userAnswer := c.readText("Tvoj odgovor: ")

		if userAnswer == "#" {
			fmt.Printf("%sVraćam te na meni...%s\n\n", utils.Yellow, utils.Reset)
			return nil
		}

		normalizedUserAnswer := normalizeAnswer(userAnswer)
		normalizedCorrectAnswer := normalizeAnswer(question.Answer)

		isCorrect := normalizedUserAnswer == normalizedCorrectAnswer

		answer := gameModels.Answer{
			WordID:        question.WordID,
			QuestionText:  question.Word,
			CorrectAnswer: question.Answer,
			UserAnswer:    userAnswer,
			IsCorrect:     isCorrect,
		}

		session.Answers = append(session.Answers, answer)

		if isCorrect {
			correctCount++
			fmt.Printf("%sOdgovor tačan%s\n\n", utils.Green, utils.Reset)
		} else {
			wrongCount++
			fmt.Printf("%sOdgovor nije tačan%s\n", utils.Red, utils.Reset)
			fmt.Printf("%sTačan odgovor je: %s%s\n\n", utils.Blue, question.Answer, utils.Reset)
		}

		if stopOnMaxWrong && wrongCount >= session.MaxWrongAnswers {
			result := gameModels.Result{
				CorrectCount:   correctCount,
				WrongCount:     wrongCount,
				TotalQuestions: len(session.Questions),
				Passed:         false,
			}

			if saveResult {
				if _, err := c.gameService.FinishLevel(session.AttemptID, session.ProgressID, result); err != nil {
					return err
				}
			}

			fmt.Println("--------------------------------------")
			fmt.Printf("%sIzgubili ste, pokušajte ponovo!%s\n", utils.Red, utils.Reset)
			fmt.Printf("Rezultat: %d tačnih, %d netačnih\n\n", correctCount, wrongCount)

			return nil
		}
	}

	result := gameModels.Result{
		CorrectCount:   correctCount,
		WrongCount:     wrongCount,
		TotalQuestions: len(session.Questions),
		Passed:         true,
	}

	if saveResult {
		if _, err := c.gameService.FinishLevel(session.AttemptID, session.ProgressID, result); err != nil {
			return err
		}
	}

	fmt.Println("--------------------------------------")

	if saveResult {
		fmt.Printf("%sUspešno savladan level!%s\n", utils.Green, utils.Reset)
	} else {
		fmt.Printf("%sPraktika završena!%s\n", utils.Green, utils.Reset)
	}

	fmt.Printf("Rezultat: %d/%d tačnih\n\n", correctCount, len(session.Questions))

	return nil
}

func (c *GameCLI) readText(label string) string {
	fmt.Print(label)

	value, _ := c.reader.ReadString('\n')
	value = strings.TrimSpace(value)

	return value
}

func (c *GameCLI) waitForEnter() {
	fmt.Println("Pritisni ENTER za povratak na meni...")
	_, _ = c.reader.ReadString('\n')
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
