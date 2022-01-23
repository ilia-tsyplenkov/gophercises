package quiz_test

import (
	"fmt"
	"io"
	"reflect"
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
			for _, userAnswer := range tc.answers {
				q := questions.NextQuiz()
				// question, correctAnswer, _ := questions.NextQuiz()
				if q.Answer != userAnswer {
					t.Errorf("expected to have %q on %q question, but got %q", userAnswer, q.Question, q.Answer)
				}
			}

		})
	}
}

func TestErrorGetQuizWhenNoMoreRecordInSlice(t *testing.T) {
	questions := quiz.SliceQuizStore{Data: [][]string{}}
	q := questions.NextQuiz()
	if q.Err != io.EOF {
		t.Fatalf("expected to have EOF, but got %s\n", q.Err)
	}

}

func TestShuffleQuizStore(t *testing.T) {
	data := [][]string{
		{"1question", "1answer"},
		{"2question", "2answer"},
		{"3question", "3answer"},
		{"4question", "4answer"},
		{"5question", "5answer"},
		{"6question", "6answer"},
		{"7question", "7answer"},
		{"8question", "8answer"},
		{"9question", "9answer"},
		{"10question", "10answer"},
	}
	questions := quiz.SliceQuizStore{Data: data}
	origin := make([][]string, len(data))
	copy(origin, data)
	questions.Shuffle()
	shuffled := questions.Data
	if reflect.DeepEqual(shuffled, origin) {
		t.Errorf("quiz data has same order after shuffilling\n")
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
				answer := answers.NextAnswer()
				got := answer.Value
				if got != want {
					t.Fatalf("expected %q, but got %q\n", want, got)
				}
			}

		})
	}

}

func TestErrorGetAnswerWhenNoMoreRecordInSlice(t *testing.T) {
	answers := quiz.SliceAnswerStore{}
	answer := answers.NextAnswer()
	if answer.Err != io.EOF {
		t.Fatalf("expected to have EOF, but got %s\n", answer.Err)
	}
}
