package quiz_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

func makeTestCsv(filename string, records [][]string) (filePath string, answers []string) {

	fd, err := os.Create(filename)
	if err == os.ErrExist {
		os.Remove(filename)
		fd, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	defer fd.Close()
	answers = make([]string, len(records))
	w := csv.NewWriter(fd)
	for i, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error preparing csv test file:", err)
		}
		answers[i] = record[1]
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return fd.Name(), answers

}
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
			makeTestCsv(testFile, tc.quizes)
			store, err := quiz.NewCsvQuizStore(testFile)
			if err != nil {
				t.Errorf("unxpected error while creating CsvQuizStore: %s\n", err)
			}
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

func TestErrorGetCsvQuizWhenNoMoreRecord(t *testing.T) {
	testFile, _ := makeTestCsv("fake.csv", nil)
	quizStore, err := quiz.NewCsvQuizStore(testFile)
	if err != nil {
		t.Errorf("unxpected error while creating CsvQuizStore: %s\n", err)
	}
	_, _, err = quizStore.NextQuiz()
	if err != io.EOF {
		t.Errorf("expecting to have EOF error while reading from empty file, but got %q\n", err)
	}
}
