package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rajkovrga/1000words-game/internal/db"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Connect(os.Getenv("DB_PATH"))
	fmt.Println("Connected!")
}
