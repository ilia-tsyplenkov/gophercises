package quiz

type FileAnswerStore struct {
}

func NewFileAnswerStore() (*FileAnswerStore, error) {
	return &FileAnswerStore{}, nil
}

func (s *FileAnswerStore) NextAnswer() (string, error) {
	return "10", nil
}
