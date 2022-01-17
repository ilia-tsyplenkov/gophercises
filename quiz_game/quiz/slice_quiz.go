package quiz

import "io"

type SliceQuiz struct {
	Data    [][2]string
	current int
}

func (q *SliceQuiz) NextQuestion() (question, answer string, err error) {
	if q.current == len(q.Data) {
		err = io.EOF
		return
	}
	question, answer = q.Data[q.current][0], q.Data[q.current][1]
	err = nil
	q.current++
	return
}

type SliceAnswers struct{}

func (a *SliceAnswers) NextAnswer() (string, error) {
	return "10", nil
}
