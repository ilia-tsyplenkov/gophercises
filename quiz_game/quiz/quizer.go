package quiz

// Quizer interface gets record from quiz data
// and returns question, correct answer and error
type QuizReader interface {
	NextQuiz() (string, string, error)
}

// AnswerReader interface returns real answer and error
type AnswerReader interface {
	NextAnswer() (string, error)
}
