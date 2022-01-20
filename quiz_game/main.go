package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

type QuizGame struct {
	quizStore   quiz.QuizReader
	answerStore quiz.AnswerReader
	// Place to show question for users
	// No questiong will shown in case of nil
	out io.Writer
}

func (g *QuizGame) CheckAnswers() (total, correct int) {
	for {
		question, rightAnswer, err := g.quizStore.NextQuiz()
		if err != nil {
			return
		}
		total++
		if g.out != nil {
			fmt.Fprintf(g.out, "%s: ", question)
		}
		answer, err := g.answerStore.NextAnswer()
		if err != nil {
			return
		}
		if rightAnswer == answer {
			correct++
		}
	}
}

var quizFile string

func init() {
	flag.StringVar(&quizFile, "quiz", "problems.csv", "csv file with question and correct answers")
}

func main() {
	flag.Parse()

	quizFd, err := os.Open(quizFile)
	if err != nil {
		log.Fatalln("error getting quiz data:", err)
	}
	defer quizFd.Close()
	quizStore := quiz.NewCsvQuizStore(quizFd)
	answerStore := quiz.NewFileAnswerStore(os.Stdin)

	game := QuizGame{quizStore, answerStore, os.Stdout}
	totalAnswers, correctAnswers := game.CheckAnswers()
	fmt.Printf("Quiz results: total - %d, correct - %d\n", totalAnswers, correctAnswers)
}
