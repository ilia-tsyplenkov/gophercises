package main

import (
	"github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"
)

type Gamer interface {
	quiz.QuizReader
	quiz.AnswerReader
	CheckAnswer(got, want interface{}) bool
}
