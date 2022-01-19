package main

import (
	"testing"
)

func TestGameQuizQuestionAnswer(t *testing.T) {
	game := QuizGame{}
	res := game.CheckAnswers("10", "10")
	if res != true {
		t.Fatalf("expected right answer but got incorrect\n")
	}
}

// func TestStartQuizGame(t *testing.T) {
// 	confirmStartFile := "confirm.txt"
// 	defer os.Remove(confirmStartFile)
// 	os.Remove(confimrStartFile)
// 	df, _ := os.Create(confirmStartFile)
//
// }

func TestStartQuizGame(t *testing.T) {
	game := QuizGame{}
	err := game.Launch()
	if err != nil {
		t.Errorf("expected success launch, but got %q", err)
	}
}
