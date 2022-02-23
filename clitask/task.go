package clitask

import (
	"errors"
	"fmt"
)

var ErrUnexpectedId = errors.New("task id must be greater than zero.")

type Task struct {
	Id   int
	Name string
	Done bool
}

func (t Task) String() string {
	return fmt.Sprintf("id: %d task: %q done: %v", t.Id, t.Name, t.Done)
}
