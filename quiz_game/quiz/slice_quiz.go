package quiz

import (
	"io"
	"time"
)

type SliceQuizStore struct {
	Data    [][]string
	current int
}

func (q *SliceQuizStore) NextQuiz() Quiz {
	next := Quiz{}
	if q.current >= len(q.Data) {
		next.Err = io.EOF
		return next
	}
	next.Question, next.Answer = q.Data[q.current][0], q.Data[q.current][1]
	q.current++
	return next
}

type SliceAnswerStore struct {
	Data    []string
	current int
}

func (a *SliceAnswerStore) NextAnswer() Answer {
	next := Answer{}
	if a.current >= len(a.Data) {
		next.Err = io.EOF
		return next
	}
	next.Value = a.Data[a.current]
	a.current++
	return next
}

type SliceDelayedAnswerStore struct {
	Store SliceAnswerStore
	Delay time.Duration
}

func (a *SliceDelayedAnswerStore) NextAnswer() Answer {
	time.Sleep(a.Delay)
	return a.Store.NextAnswer()
}
