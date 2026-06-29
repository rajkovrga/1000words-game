package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Question struct {
	Word   string
	Answer string
}

func main() {

	questions := []Question{
		{Word: "cat", Answer: "mačka"},
		{Word: "dog", Answer: "pas"},
		{Word: "house", Answer: "kuća"},
		{Word: "car", Answer: "auto"},
	}

	fmt.Printf("Dobrodošli na kviz znanja \n")
	fmt.Printf("---------------------------- \n")
	fmt.Printf("Trenutno ste na nivou: 1 \n")
	fmt.Printf("---------------------------- \n")

	reader := bufio.NewReader(os.Stdin)

	randomOrder := rand.Perm(len(questions))

	totalQuestions := len(randomOrder)
	bugs := 0
	for i := 0; i < totalQuestions; i++ {
		randomIndex := randomOrder[i]
		question := questions[randomIndex]

		fmt.Printf("Reč: %s, napisi odgovor:", question.Word)

		userAnswer, _ := reader.ReadString('\n')
		userAnswer = strings.TrimSpace(userAnswer)
		userAnswer = strings.ToLower(userAnswer)

		if strings.ToLower(question.Answer) == userAnswer {
			fmt.Printf("\033[32mOdgovor tačan\033[0m\n\n")
		} else {
			fmt.Printf("\033[31mOdgovor nije tačan\033[0m\n\n")
			fmt.Printf("\033[34mTačan odgovor je: %s\033[0m\n\n", question.Answer)
			bugs += 1
		}

		if bugs == 3 {
			fmt.Printf("\033[31mIzgubili ste, pokušajte ponovo!\033[0m\n\n")
			break
		}

	}

	fmt.Printf("\033[32mUspešno savladan nivo 1\033[0m\n\n")
}
