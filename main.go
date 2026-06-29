package main

import (
	"1000words-game/data"
	"1000words-game/game"
)

func main() {
	game.StartLevelCommand(1, data.GetLevelWords(1))
}
