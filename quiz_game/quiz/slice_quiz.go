package quiz

import "io"

type SliceQuiz struct {
	Data    [][2]string
	current int
}

func (q *SliceQuiz) NextQuestion() (question, answer string, err error) {
	if q.current >= len(q.Data) {
		err = io.EOF
		return
	}
	question, answer = q.Data[q.current][0], q.Data[q.current][1]
	q.current++
	return
}

type SliceAnswers struct {
	Data    []string
	current int
}

func (a *SliceAnswers) NextAnswer() (answer string, err error) {
	if a.current >= len(a.Data) {
		err = io.EOF
		return
	}
	answer = a.Data[a.current]
	a.current++
	return
}
