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
		quizData    [][]string
		userAnswers []string
		total       int
		correct     int
	}{
		{
			quizData:    [][]string{{"10 + 10", "20"}},
			userAnswers: []string{"20"},
			total:       1,
			correct:     1,
		},
		{
			quizData:    [][]string{{"10 + 10", "20"}},
			userAnswers: []string{"25"},
			total:       1,
			correct:     0,
		},
		{
			quizData:    [][]string{{"10 + 10", "20"}, {"20-5", "15"}},
			userAnswers: []string{"20", "15"},
			total:       2,
			correct:     2,
		},
		{
			quizData:    [][]string{{"10 + 10", "20"}, {"20-5", "15"}},
			userAnswers: []string{"20", "10"},
			total:       2,
			correct:     1,
		},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("%dQuestions%dTotalAnswers%dCorrectAnswers", len(tc.quizData), tc.total, tc.correct)
		t.Run(testName, func(t *testing.T) {

			quizStore := &quiz.SliceQuizStore{Data: tc.quizData}
			answerStore := &quiz.SliceAnswerStore{Data: tc.userAnswers}
			game := QuizGame{quizStore, answerStore, nil, 0, nil}
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
	game := QuizGame{quizStore: quizStore, answerStore: answerStore, out: nil, timeout: gameTimeout}
	total, _ := game.CheckAnswers()

	want := 0
	if total != want {
		t.Fatalf("expected to have %d answered questions, but got - %d\n", total, want)
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
	game.Greeting()
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
