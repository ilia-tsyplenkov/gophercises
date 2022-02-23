package clitask

import "io"

type MemStore struct {
	Todo []Task
	Done []Task
}

func NewMemStore() *MemStore {
	return &MemStore{}
}

func (s *MemStore) AllTasks() []Task {
	return s.Todo
}

func (s *MemStore) ToDo() ([]Task, error) {
	return s.Todo, nil
}

func (s *MemStore) Add(task string) error {
	s.Todo = append(s.Todo, Task{Name: task})
	return nil
}

func (s *MemStore) Completed() []Task {
	return s.Done
}

func (s *MemStore) Do(id int) (Task, error) {
	if id < 1 {
		return Task{}, ErrUnexpectedId
	}
	if id > len(s.Todo) {
		return Task{}, io.ErrUnexpectedEOF
	}
	id--
	task := s.Todo[id]
	s.Todo = append(s.Todo[:id], s.Todo[id+1:]...)
	task.Done = true
	s.Done = append(s.Done, task)
	return task, nil
}
