package quiz_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

func TestGetQuestionAndAnswerFromSliceQuizStore(t *testing.T) {
	testCases := []struct {
		question_data [][]string
		answers       []string
	}{
		{
			[][]string{{"10 + 10", "20"}},
			[]string{"20"},
		},
		{
			[][]string{{"10 + 10", "20"}, {"10 - 5", "5"}},
			[]string{"20", "5"},
		},
		{
			[][]string{{"10 + 10", "20"}, {"10 - 5", "5"}, {"10 + 5", "15"}},
			[]string{"20", "5", "15"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%dQuestions%dAnswers", len(tc.question_data), len(tc.answers)), func(t *testing.T) {
			questions := quiz.SliceQuizStore{Data: tc.question_data}
			for _, answer := range tc.answers {
				question, correctAnswer, _ := questions.NextQuiz()
				if answer != correctAnswer {
					t.Errorf("expected to have %q on %q question, but got %q", answer, question, correctAnswer)
				}
			}

		})
	}
}

func TestErrorGetQuizWhenNoMoreRecordInSlice(t *testing.T) {
	questions := quiz.SliceQuizStore{Data: [][]string{}}
	_, _, err := questions.NextQuiz()
	if err != io.EOF {
		t.Fatalf("expected to have EOF, but got %s\n", err)
	}

}

func TestGetAnswerFromSliceAnswerStore(t *testing.T) {
	testCases := []struct {
		answers         []string
		expectedAnswers []string
	}{
		{answers: []string{"10", "20"},
			expectedAnswers: []string{"10", "20"}},

		{answers: []string{"10", "20", "30"},
			expectedAnswers: []string{"10", "20", "30"}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%dExistingAnswers%dExpedted", len(tc.answers), len(tc.expectedAnswers)), func(t *testing.T) {

			answers := quiz.SliceAnswerStore{Data: tc.answers}
			for _, want := range tc.expectedAnswers {
				got, _ := answers.NextAnswer()
				if got != want {
					t.Fatalf("expected %q, but got %q\n", want, got)
				}
			}

		})
	}

}

func TestErrorGetAnswerWhenNoMoreRecordInSlice(t *testing.T) {
	answers := quiz.SliceAnswerStore{}
	_, err := answers.NextAnswer()
	if err != io.EOF {
		t.Fatalf("expected to have EOF, but got %s\n", err)
	}
}
