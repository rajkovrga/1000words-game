package cli

import "fmt"

func (a *App) Run() error {
	for {
		fmt.Println("======================================")
		fmt.Println("1000words game")
		fmt.Println("======================================")
		fmt.Println("1. Praktika")
		fmt.Println("2. Igra")
		fmt.Println("0. Izlaz")
		fmt.Println("--------------------------------------")

		choice := a.readText("Izaberi opciju: ")

		switch choice {
		case "1":
			if err := a.runPractice(); err != nil {
				printError(fmt.Sprintf("Greška: %s", err.Error()))
				fmt.Println()
				a.waitForEnter()
			}

		case "2":
			if err := a.runGame(); err != nil {
				printError(fmt.Sprintf("Greška: %s", err.Error()))
				fmt.Println()
				a.waitForEnter()
			}

		case "0":
			fmt.Println("Izlaz iz aplikacije.")
			return nil

		default:
			printError("Nepoznata opcija.")
			fmt.Println()
		}
	}
}
