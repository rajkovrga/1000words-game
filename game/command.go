package game

import (
	models "1000words-game/models/game"
	"1000words-game/utils"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func StartLevelCommand(level int, questions []models.Question) {

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
			fmt.Printf("%sOdgovor tačan%s\n\n", utils.Green, utils.Reset)
		} else {
			fmt.Printf("%sOdgovor nije tačan%s\n\n", utils.Red, utils.Reset)
			fmt.Printf("%sTačan odgovor je: %s%s\n\n", utils.Blue, question.Answer, utils.Reset)
			bugs++
		}

		if bugs == 3 {
			fmt.Printf("%sIzgubili ste, pokušajte ponovo!%s\n\n", utils.Red, utils.Reset)
			break
		}

	}

	fmt.Printf("%sSa %s greške, uspešno savladan nivo 1%s\n\n", utils.Green, bugs, utils.Reset)
}
