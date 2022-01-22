package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

type QuizGame struct {
	quizStore   quiz.QuizReader
	answerStore quiz.AnswerReader
	// Place to show question for users
	// No questiong will shown in case of nil
	out     io.Writer
	timeout time.Duration
	// reader to define user readiness
	in io.Reader
}

func (g *QuizGame) CheckAnswers() (total, correct int) {
	var stop <-chan time.Time
	if g.timeout > 0 {
		stop = time.After(g.timeout)
	}
	answers := make(chan quiz.Answer)
	go func() {
		for {
			ans := g.answerStore.NextAnswer()
			answers <- ans
			if ans.Err != nil {
				break
			}
		}
	}()
	for {
		question := g.quizStore.NextQuiz()
		if question.Err != nil {
			return
		}
		if g.out != nil {
			fmt.Fprintf(g.out, "%s: ", question.Question)
		}
		select {
		case <-stop:
			return
		case userAnswer := <-answers:
			if userAnswer.Err != nil {
				return
			}
			total++
			if question.Answer == userAnswer.Value {
				correct++
			}
		}
	}
}

func (g *QuizGame) Greeting() {
	fmt.Fprint(g.out, "Welcome to the Quiz Game. Press any key to start:")
}

func (g *QuizGame) waitUserReadiness() {
	buffer := bufio.NewReader(g.in)
	for {
		_, err := buffer.ReadString('\n')
		if err == nil {
			break
		}
	}
}

var quizFile string
var timeout time.Duration

func init() {
	flag.StringVar(&quizFile, "quiz", "problems.csv", "csv file with question and correct answers")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "quiz timeout")
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

	game := QuizGame{quizStore, answerStore, os.Stdout, timeout, os.Stdin}
	game.Greeting()
	game.waitUserReadiness()
	totalAnswers, correctAnswers := game.CheckAnswers()
	fmt.Printf("Quiz results: total - %d, correct - %d\n", totalAnswers, correctAnswers)
}
