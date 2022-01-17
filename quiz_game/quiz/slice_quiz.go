package quiz

type SliceQuiz struct {
	Data    [][2]string
	current int
}

func (q *SliceQuiz) NextQuestion() (question, answer string, err error) {
	question, answer = q.Data[q.current][0], q.Data[q.current][1]
	err = nil
	q.current++
	return
}
