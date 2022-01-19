package quiz_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
	sugar "github.com/ilia-tsyplenkov/gophercises/quiz_game/test_sugar"
)

func TestGetRecordFromFileAnswerStore(t *testing.T) {

	testCases := [][]string{
		{"10"},
		{"10", "20", "30"},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%dAnswers", len(tc))
		t.Run(testName, func(t *testing.T) {
			answersFile := testName + ".txt"
			err := sugar.MakeTestAnswerFile(answersFile, tc)
			if err != nil {
				t.Fatalf("error creating test answers file: %s\n", err)
			}
			answerStore, _ := quiz.NewFileAnswerStore(answersFile)
			for _, want := range tc {
				got, err := answerStore.NextAnswer()
				if err != nil {
					t.Fatalf("error getting answer: %s\n", err)
				}
				if got != want {
					t.Errorf("expected %q, but got %q\n", want, got)
				}

			}
		})
	}

}

func TestErrorGetAnswerWhenNoRecordFile(t *testing.T) {
	answerStore, err := quiz.NewFileAnswerStore("")
	if err != nil {
		t.Fatalf("error creating FileAnswerStore: %s\n", err)
	}
	_, err = answerStore.NextAnswer()
	if err != io.EOF {
		t.Fatalf("expecting EOF while reading empty file, but got: %s\n", err)
	}
}
