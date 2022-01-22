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
	quiz.QuizReader
	quiz.AnswerReader
	// reader to define user readiness
	// won't expect readiness in case of nil
	in io.Reader
	// Place to show question for users
	// No questiong will shown in case of nil
	out io.Writer
	// Stop quiz timeout
	timeout time.Duration
}

func (g *QuizGame) CheckAnswers() (total, correct int) {
	total = g.Total()
	var stop <-chan time.Time
	if g.timeout > 0 {
		stop = time.After(g.timeout)
	}
	answers := make(chan quiz.Answer)
	go func() {
		for {
			ans := g.NextAnswer()
			answers <- ans
			if ans.Err != nil {
				break
			}
		}
	}()
	for {
		question := g.NextQuiz()
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
	if g.in == nil {
		return
	}
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
	quizStore, _ := quiz.NewSliceQuizFromCsv(quizFd)
	answerStore := quiz.NewFileAnswerStore(os.Stdin)

	game := QuizGame{
		QuizReader:   quizStore,
		AnswerReader: answerStore,
		in:           os.Stdin,
		out:          os.Stdout,
		timeout:      timeout,
	}
	game.Greeting()
	game.waitUserReadiness()
	totalAnswers, correctAnswers := game.CheckAnswers()
	fmt.Printf("\nQuiz results: total - %d, correct - %d\n", totalAnswers, correctAnswers)
}
