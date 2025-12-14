package main

import (
	"fmt"

	"github.com/rajkovrga/1000words-game/internal/db"
)

func main() {
	db.Connect("./../../database.db")
	fmt.Println("Connected!")
}
