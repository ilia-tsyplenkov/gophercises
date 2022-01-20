package quiz_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
	"github.com/ilia-tsyplenkov/gophercises/quiz_game/test_sugar"
	sugar "github.com/ilia-tsyplenkov/gophercises/quiz_game/test_sugar"
)

func TestGetQuestionAndAnswerFromCsvQuizStore(t *testing.T) {
	testCases := []struct {
		quizes  [][]string
		answers []string
	}{
		{
			[][]string{
				{"10+10", "20"},
				{"10-5", "5"},
				{"15+10", "25"},
			},
			[]string{"20", "5", "25"},
		},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("%dQuestions%dAnswers", len(tc.quizes), len(tc.answers))
		t.Run(testName, func(t *testing.T) {
			testFile := "testName" + ".csv"
			sugar.MakeTestCsv(testFile, tc.quizes)
			f, _ := os.Open(testFile)
			defer func() {
				defer f.Close()
				os.Remove(testFile)
			}()
			store := quiz.NewCsvQuizStore(f)
			for _, want := range tc.answers {
				_, got, err := store.NextQuiz()
				if err != nil {
					t.Fatalf("unexpected error during csv quiz reading: %s\n", err)
				}
				if got != want {
					t.Errorf("expected %q but got %q\n", want, got)
				}

			}

		})
	}

}

func TestErrorGetQuizWhenNoMoreRecordInCsv(t *testing.T) {
	testFile, _ := test_sugar.MakeTestCsv("fake.csv", nil)
	f, _ := os.Open(testFile)
	defer func() {
		defer f.Close()
		os.Remove(testFile)
	}()
	quizStore := quiz.NewCsvQuizStore(f)
	_, _, err := quizStore.NextQuiz()
	if err != io.EOF {
		t.Errorf("expecting to have EOF error while reading from empty file, but got %q\n", err)
	}
}
