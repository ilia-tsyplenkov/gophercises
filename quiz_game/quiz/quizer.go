package quiz

// Quizer interface gets record from quiz data
// and returns question, correct answer and error
type QuizReader interface {
	NextQuiz() Quiz
}

// AnswerReader interface returns real answer and error
type AnswerReader interface {
	NextAnswer() Answer
}

type Quiz struct {
	Question string
	Answer   string
	Err      error
}

type Answer struct {
	Value string
	Err   error
}
