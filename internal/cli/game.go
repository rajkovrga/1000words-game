package cli

import (
	"fmt"
	"strconv"

	dbModels "1000words-game/models/db"
	gameModels "1000words-game/models/game"
	"1000words-game/utils"
)

func (a *App) startGameForUser(user *dbModels.User) error {
	fmt.Println()
	fmt.Printf("%sKorisnik: %s%s\n", utils.Green, user.Email, utils.Reset)

	progress, err := a.chooseOrCreateProgress(user.ID)
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

	session, err := a.gameService.StartLevel(
		user.ID,
		*progress,
	)
	if err != nil {
		return err
	}

	if err := a.runQuestions(session, true, true); err != nil {
		return err
	}

	a.waitForEnter()
	return nil
}

func (a *App) chooseOrCreateProgress(userID int) (*gameModels.ProgressOption, error) {
	options, err := a.progressService.GetUserProgressOptions(userID)
	if err != nil {
		return nil, err
	}

	if len(options) == 0 {
		fmt.Println()
		fmt.Println("Nemaš podešen jezik za igru.")
		fmt.Println("Prvo izaberi target i native jezik.")

		targetCode, nativeCode, err := a.chooseLanguages()
		if err != nil {
			return nil, err
		}

		return a.progressService.CreateProgress(
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
		value := a.readText("Opcija: ")

		selectedIndex, err := strconv.Atoi(value)
		if err != nil {
			printError("Unesi broj opcije.")
			continue
		}

		if selectedIndex < 1 || selectedIndex > len(options) {
			printError("Ta opcija ne postoji.")
			continue
		}

		return &options[selectedIndex-1], nil
	}
}

func (a *App) runQuestions(
	session *gameModels.Session,
	saveResult bool,
	stopOnMaxWrong bool,
) error {
	correctCount := 0
	wrongCount := 0

	fmt.Println()
	printHeader("KVIZ POČINJE")
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

		userAnswer := a.readText("Tvoj odgovor: ")

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
				if _, err := a.gameService.FinishLevel(session.AttemptID, session.ProgressID, result); err != nil {
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

	var nextLevel *dbModels.Level

	if saveResult {
		finishedNextLevel, err := a.gameService.FinishLevel(session.AttemptID, session.ProgressID, result)
		if err != nil {
			return err
		}

		nextLevel = finishedNextLevel
	}

	fmt.Println("--------------------------------------")

	if saveResult {
		printSuccess("Čestitamo! Položili ste level.")

		if nextLevel != nil {
			fmt.Printf(
				"%sPrebačeni ste na sledeći nivo: Level %d - %s%s\n",
				utils.Blue,
				nextLevel.LevelNumber,
				nextLevel.Name,
				utils.Reset,
			)
		} else {
			printSuccess("Savladali ste poslednji level. Bravo!")
		}
	} else {
		printSuccess("Praktika završena!")
	}

	fmt.Printf("Rezultat: %d/%d tačnih\n\n", correctCount, len(session.Questions))

	return nil
}