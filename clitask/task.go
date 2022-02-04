package clitask

func GetTask() string {
	return "write tests"
}

type MemStore struct {
	tasks []string
}

func (s *MemStore) TaskList() []string {
	return s.tasks
}
