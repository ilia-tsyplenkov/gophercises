package quiz_test

import (
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

func TestGetRecordFromFileAnswerStore(t *testing.T) {

	answerStore, _ := quiz.NewFileAnswerStore()
	got, _ := answerStore.NextAnswer()
	want := "10"
	if got != want {
		t.Errorf("expected to have %q, but got %q\n", want, got)
	}

}
