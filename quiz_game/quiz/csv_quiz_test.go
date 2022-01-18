package quiz_test

import (
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

func TestGetQuestionAndAnswerFromCsvQuizStore(t *testing.T) {
	store := quiz.NewCsvQuizStore()
	_, got, _ := store.NextQuiz()
	want := "10"
	if got != want {
		t.Fatalf("expected %q but got %q\n", want, got)
	}

}
