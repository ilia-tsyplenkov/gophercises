package quiz_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

func TestGetQuestionAndAnswerFromSliceQuiz(t *testing.T) {
	testCases := []struct {
		question_data [][2]string
		answers       []string
	}{
		{
			[][2]string{{"10 + 10", "20"}},
			[]string{"20"},
		},
		{
			[][2]string{{"10 + 10", "20"}, {"10 - 5", "5"}},
			[]string{"20", "5"},
		},
		{
			[][2]string{{"10 + 10", "20"}, {"10 - 5", "5"}, {"10 + 5", "15"}},
			[]string{"20", "5", "15"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%dQuestions%dAnswers", len(tc.question_data), len(tc.answers)), func(t *testing.T) {
			questions := quiz.SliceQuiz{Data: tc.question_data}
			for _, answer := range tc.answers {
				question, correctAnswer, _ := questions.NextQuestion()
				if answer != correctAnswer {
					t.Errorf("expected to have %q on %q question, but got %q", answer, question, correctAnswer)
				}
			}

		})
	}
}

func TestErrorThanNoMoreQuestions(t *testing.T) {
	questions := quiz.SliceQuiz{Data: [][2]string{}}
	_, _, err := questions.NextQuestion()
	if err != io.EOF {
		t.Fatalf("expected to have EOF, but got %s\n", err)
	}

}

func TestGetAnswerFromSliceAnswers(t *testing.T) {
	answers := quiz.SliceAnswers{}
	got, _ := answers.NextAnswer()
	want := "10"
	if got != want {
		t.Fatalf("expected %q, but got %q\n", want, got)
	}

}

// func TestFistSliceQuestionHasRightFromSlice(t *testing.T) {
// 	questions := quiz.SliceQuiz{Data: [][2]string{{"10 + 10", "20"}}}
// 	// answers := quiz.SliceAnswer{Data: []string{"20"}}
//
//
// }
