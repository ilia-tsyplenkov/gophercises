package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
	shuffle bool
}

func (g *QuizGame) checkAnswers() (total, correct int) {
	total = g.Total()
	if g.shuffle {
		g.QuizReader.Shuffle()
	}
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
			if correctIt(question.Answer) == correctIt(userAnswer.Value) {
				correct++
			}
		}
	}
}

func (g *QuizGame) greeting() {
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

func (g *QuizGame) Start() (total, correct int) {
	g.greeting()
	g.waitUserReadiness()
	return g.checkAnswers()
}

func correctIt(s string) string {
	s = strings.Trim(s, " \n\t")
	return strings.ToLower(s)
}
