package cli

import (
	"fmt"

	dbModels "1000words-game/models/db"
)

func (a *App) runGame() error {
	if a.currentUser != nil {
		return a.startGameForUser(a.currentUser)
	}

	for {
		fmt.Println()
		printHeader("IGRA")
		fmt.Println("1. Prijava")
		fmt.Println("2. Registracija")
		fmt.Println("0. Nazad")
		fmt.Println("--------------------------------------")

		choice := a.readText("Izaberi opciju: ")

		switch choice {
		case "1":
			user, err := a.login()
			if err != nil {
				return err
			}

			a.currentUser = user
			return a.startGameForUser(a.currentUser)

		case "2":
			user, err := a.register()
			if err != nil {
				return err
			}

			a.currentUser = user
			return a.startGameForUser(a.currentUser)

		case "0":
			return nil

		default:
			printError("Nepoznata opcija.")
		}
	}
}

func (a *App) login() (*dbModels.User, error) {
	fmt.Println()
	printHeader("PRIJAVA")

	email := a.readText("Email: ")
	password := a.readText("Password: ")

	user, err := a.userService.Login(email, password)
	if err != nil {
		return nil, err
	}

	printSuccess("Uspešno ste se prijavili.")
	return user, nil
}

func (a *App) register() (*dbModels.User, error) {
	fmt.Println()
	printHeader("REGISTRACIJA")

	email := a.readText("Email: ")
	password := a.readText("Password: ")

	user, err := a.userService.Register(email, password)
	if err != nil {
		return nil, err
	}

	printSuccess("Uspešno ste se registrovali.")
	return user, nil
}
