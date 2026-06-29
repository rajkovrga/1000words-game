package cli

import (
	"fmt"

	"1000words-game/utils"
)

func printHeader(title string) {
	fmt.Println("======================================")
	fmt.Println(title)
	fmt.Println("======================================")
}

func printError(message string) {
	fmt.Printf("%s%s%s\n", utils.Red, message, utils.Reset)
}

func printSuccess(message string) {
	fmt.Printf("%s%s%s\n", utils.Green, message, utils.Reset)
}

func printInfo(message string) {
	fmt.Printf("%s%s%s\n", utils.Blue, message, utils.Reset)
}

func printWarning(message string) {
	fmt.Printf("%s%s%s\n", utils.Yellow, message, utils.Reset)
}
