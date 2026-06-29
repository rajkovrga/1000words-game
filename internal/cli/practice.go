package cli

import "fmt"

func (a *App) runPractice() error {
	fmt.Println()
	printHeader("PRAKTIKA")

	levelNumber, err := a.chooseLevel()
	if err != nil {
		return err
	}

	targetCode, nativeCode, err := a.chooseLanguages()
	if err != nil {
		return err
	}

	session, err := a.gameService.StartPractice(
		levelNumber,
		targetCode,
		nativeCode,
	)
	if err != nil {
		return err
	}

	if err := a.runQuestions(session, false, false); err != nil {
		return err
	}

	a.waitForEnter()
	return nil
}
