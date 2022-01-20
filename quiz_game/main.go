package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

type QuizGame struct {
	quizStore   quiz.QuizReader
	answerStore quiz.AnswerReader
}

func (g *QuizGame) CheckAnswers() (total, correct int) {
	for {
		question, rightAnswer, err := g.quizStore.NextQuiz()
		if err != nil {
			return
		}
		total++
		fmt.Fprintf(os.Stdout, "%s: ", question)
		answer, err := g.answerStore.NextAnswer()
		if err != nil {
			return
		}
		if rightAnswer == answer {
			correct++
		}
	}
}

func (g *QuizGame) Launch() error {
	return nil
}

var defaultFile string = "problems.csv"

func main() {

	quizStore, err := quiz.NewCsvQuizStore(defaultFile)
	if err != nil {
		log.Fatalln("error getting quiz data:", err)
	}
	answerStore, err := quiz.NewFileAnswerStore("")
	if err != nil {
		log.Fatalln("error getting file(file descriptor) for answers: ", err)
	}
	game := QuizGame{quizStore, answerStore}
	totalAnswers, correctAnswers := game.CheckAnswers()
	fmt.Printf("Quiz results: total - %d, correct -%d\n", totalAnswers, correctAnswers)
}
