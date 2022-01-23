package quiz

import (
	"encoding/csv"
	"io"
	"math/rand"
	"os"
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

func (q *SliceQuizStore) Total() int {
	return len(q.Data)
}

func (q *SliceQuizStore) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.Data), func(i, j int) { q.Data[i], q.Data[j] = q.Data[j], q.Data[i] })
}

// Factory function which get file descriptor and reads all records
// into SliceQuizStore
func NewSliceQuizFromCsv(fd *os.File) (*SliceQuizStore, error) {
	data, err := csv.NewReader(fd).ReadAll()
	if err != nil {
		return nil, err
	}
	return &SliceQuizStore{Data: data}, nil

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
	SliceAnswerStore
	Delay time.Duration
}

func (a *SliceDelayedAnswerStore) NextAnswer() Answer {
	time.Sleep(a.Delay)
	return a.SliceAnswerStore.NextAnswer()
}
