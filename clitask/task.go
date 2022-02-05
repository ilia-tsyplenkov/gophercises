package clitask

import "fmt"

func GetTask() string {
	return "write tests"
}

type Task struct {
	name string
	done bool
}

func (t Task) String() string {
	return fmt.Sprintf("task: %q done: %v", t.name, t.done)
}

type MemStore struct {
	todo []Task
	done []Task
}

func (s *MemStore) AllTasks() []Task {
	return s.todo
}

func (s *MemStore) ToDo() []Task {
	return s.todo
}

func (s *MemStore) Add(task string) {
	s.todo = append(s.todo, Task{name: task})
}

func (s *MemStore) Completed() []Task {
	return s.done
}

func (s *MemStore) Do(id int) {
	id--
	task := s.todo[id]
	s.todo = append(s.todo[:id], s.todo[id+1:]...)
	task.done = true
	s.done = append(s.done, task)
}
