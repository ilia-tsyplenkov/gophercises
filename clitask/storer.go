package clitask

type Storer interface {
	ToDo() ([]Task, error)
	Add(string) error
	Do(int) (Task, error)
}
