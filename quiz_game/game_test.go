package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

func TestGameQuizQuestionsCorrectAnswers(t *testing.T) {
	testCases := []struct {
		title       string
		quizData    [][]string
		userAnswers []string
		total       int
		correct     int
	}{
		{
			title:       "1Question1TotalAnswer1CorrectAnswer",
			quizData:    [][]string{{"10 + 10", "20"}},
			userAnswers: []string{"20"},
			total:       1,
			correct:     1,
		},
		{
			title:       "1Question1TotalAnswer0CorrectAnswer",
			quizData:    [][]string{{"10 + 10", "20"}},
			userAnswers: []string{"25"},
			total:       1,
			correct:     0,
		},
		{
			title:       "2Questions2TotalAnswers2CorrectAnswers",
			quizData:    [][]string{{"10 + 10", "20"}, {"20-5", "15"}},
			userAnswers: []string{"20", "15"},
			total:       2,
			correct:     2,
		},
		{
			title:       "2Questions2TotalAnswers1CorrectAnswer",
			quizData:    [][]string{{"10 + 10", "20"}, {"20-5", "15"}},
			userAnswers: []string{"20", "10"},
			total:       2,
			correct:     1,
		},
		{
			title:       "CorrectAnswersBeforeCompare",
			quizData:    [][]string{{"What's you favorite program language", " Go\t"}, {"What's our planet name?", "\nEARTH "}},
			userAnswers: []string{"  gO\n", "\t  eaRtH \n"},
			total:       2,
			correct:     2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {

			quizStore := &quiz.SliceQuizStore{Data: tc.quizData}
			answerStore := &quiz.SliceAnswerStore{Data: tc.userAnswers}
			game := QuizGame{QuizReader: quizStore, AnswerReader: answerStore, out: nil, timeout: 0, in: nil}
			total, correct := game.CheckAnswers()
			if total != tc.total {
				t.Fatalf("expected to have %d total answered questions, but got %d\n", tc.total, total)
			}
			if correct != tc.correct {
				t.Fatalf("expected to have %d correct answers but got %d\n", tc.correct, correct)
			}
		})
	}
}

func TestQuizGameTimeIsUp(t *testing.T) {
	gameTimeout := 10 * time.Millisecond
	quizStore := &quiz.SliceQuizStore{Data: [][]string{{"10 + 10", "20"}}}
	answerStore := &quiz.SliceDelayedAnswerStore{Store: quiz.SliceAnswerStore{Data: []string{"20"}}, Delay: 2 * gameTimeout}
	game := QuizGame{quizStore, answerStore, nil, nil, gameTimeout}
	_, answered := game.CheckAnswers()

	want := 0
	if answered != want {
		t.Fatalf("expected to have %d answered questions, but got - %d\n", answered, want)
	}

}

func TestGameShowGreeting(t *testing.T) {
	greetingFile := "greeting.txt"
	fd, _ := os.Create(greetingFile)
	defer func() {
		os.Remove(greetingFile)
		fd.Close()
	}()
	game := QuizGame{
		out: fd,
	}
	// Write greeting message
	game.greeting()
	// Back in the begininng of the file
	fd.Seek(0, 0)

	// Check that greeting message has been written
	buffer := bufio.NewReader(fd)
	userReadiness, err := buffer.ReadString('\n')

	if err != nil && userReadiness == "" {
		t.Fatal("expect to have some greeting, but got nothing.")
	}

}

func TestGameAcceptUserReadiness(t *testing.T) {
	fileName := "rediness.txt"
	fd, _ := os.Create(fileName)
	defer func() {
		fd.Close()
		os.Remove(fileName)
	}()
	fmt.Fprintf(fd, "Y\n")
	// Back in the begininng of the file
	fd.Seek(0, 0)
	game := QuizGame{
		in: fd,
	}

	ready := make(chan struct{})
	go func() {
		game.waitUserReadiness()
		ready <- struct{}{}
	}()
	for {
		select {
		case <-ready:
			return
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("expected accepting user readiness, but user answer has been ignored.\n")
		}
	}
}

func TestGameNoUserReadinessProvided(t *testing.T) {
	game := QuizGame{}
	game.waitUserReadiness()
}

func TestCorrectIt(t *testing.T) {
	testCases := []struct {
		initial string
		want    string
	}{
		{" 10 ", "10"},
		{"  Spaces \n", "spaces"},
		{"\tTabs\t", "tabs"},
		{"\t  SpacesAndTabs\t  ", "spacesandtabs"},
		{"\nLineFeed\n", "linefeed"},
	}
	for _, tc := range testCases {
		got := CorrectIt(tc.initial)
		want := tc.want
		if got != want {
			t.Errorf("expected to have %q after correction of %q but got %q\n", want, tc.initial, got)
		}
	}
}
