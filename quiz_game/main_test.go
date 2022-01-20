package main

import (
	"fmt"
	"testing"

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
			game := QuizGame{quizStore, answerStore}
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

func TestStartQuizGame(t *testing.T) {
	game := QuizGame{}
	err := game.Launch()
	if err != nil {
		t.Errorf("expected success launch, but got %q", err)
	}
}
