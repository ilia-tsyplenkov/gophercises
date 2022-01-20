package main

import "github.com/ilia-tsyplenkov/gophercises/quiz_game/quiz"

type QuizGame struct {
	quizStore   quiz.QuizReader
	answerStore quiz.AnswerReader
}

func (g *QuizGame) CheckAnswers() (total, correct int) {
	for {
		_, rightAnswer, err := g.quizStore.NextQuiz()
		if err != nil {
			return
		}
		total++
		answer, err := g.answerStore.NextAnswer()
		if err != nil {
			return
		}
		if rightAnswer == answer {
			correct++
		}
	}
}

func (g *QuizGame) Launch() error {
	return nil
}

func main() {

}
