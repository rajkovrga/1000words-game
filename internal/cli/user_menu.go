package cli

import (
	"errors"
	"fmt"

	"1000words-game/utils"
)

func (a *App) runUserMenu() error {
	if a.currentUser == nil {
		return errors.New("korisnik nije prijavljen")
	}

	for {
		fmt.Println()
		printHeader("KORISNIČKI MENI")
		fmt.Printf("%sKorisnik: %s%s\n", utils.Green, a.currentUser.Email, utils.Reset)
		fmt.Println()
		fmt.Println("1. Nastavi igru")
		fmt.Println("2. Moj napredak")
		fmt.Println("3. Dodaj novi jezik za učenje")
		fmt.Println("4. Odjava")
		fmt.Println("0. Nazad")
		fmt.Println("--------------------------------------")

		choice := a.readText("Izaberi opciju: ")

		switch choice {
		case "1":
			if err := a.startGameForUser(a.currentUser); err != nil {
				return err
			}

		case "2":
			if err := a.showUserProgress(a.currentUser.ID); err != nil {
				return err
			}

			a.waitForEnter()

		case "3":
			if err := a.addNewProgress(a.currentUser.ID); err != nil {
				return err
			}

			a.waitForEnter()

		case "4":
			a.currentUser = nil
			printSuccess("Uspešno ste se odjavili.")
			return nil

		case "0":
			return nil

		default:
			printError("Nepoznata opcija.")
		}
	}
}

func (a *App) showUserProgress(userID int) error {
	options, err := a.progressService.GetUserProgressOptions(userID)
	if err != nil {
		return err
	}

	fmt.Println()
	printHeader("MOJ NAPREDAK")

	if len(options) == 0 {
		printWarning("Još nemaš podešen nijedan jezik za igru.")
		return nil
	}

	for index, option := range options {
		fmt.Printf(
			"%d. %s -> %s | Trenutni level: %d\n",
			index+1,
			option.TargetName,
			option.NativeName,
			option.LevelNumber,
		)
	}

	return nil
}

func (a *App) addNewProgress(userID int) error {
	fmt.Println()
	printHeader("DODAJ NOVI JEZIK ZA UČENJE")

	targetCode, nativeCode, err := a.chooseLanguages()
	if err != nil {
		return err
	}

	progress, err := a.progressService.CreateProgress(
		userID,
		targetCode,
		nativeCode,
	)
	if err != nil {
		return err
	}

	printSuccess("Novi jezik za učenje je dodat.")

	fmt.Printf(
		"%s%s -> %s | Level %d%s\n",
		utils.Blue,
		progress.TargetName,
		progress.NativeName,
		progress.LevelNumber,
		utils.Reset,
	)

	return nil
}
