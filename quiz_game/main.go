package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

var quizFile string
var timeout time.Duration
var shuffle bool

func init() {
	flag.StringVar(&quizFile, "quiz", "problems.csv", "csv file with question and correct answers")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "quiz timeout")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle quiz questions")
}

func main() {
	flag.Parse()

	quizFd, err := os.Open(quizFile)
	if err != nil {
		log.Fatalln("error getting quiz data:", err)
	}
	defer quizFd.Close()
	quizStore, _ := quiz.NewSliceQuizFromCsv(quizFd)
	answerStore := quiz.NewFileAnswerStore(os.Stdin)

	game := QuizGame{
		QuizReader:   quizStore,
		AnswerReader: answerStore,
		in:           os.Stdin,
		out:          os.Stdout,
		timeout:      timeout,
		shuffle:      shuffle,
	}
	game.greeting()
	game.waitUserReadiness()
	totalAnswers, correctAnswers := game.CheckAnswers()
	fmt.Printf("\nQuiz results: total - %d, correct - %d\n", totalAnswers, correctAnswers)
}
