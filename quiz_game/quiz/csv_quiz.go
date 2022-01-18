package quiz

type CsvQuizStore struct {
}

func (s *CsvQuizStore) NextQuiz() (string, string, error) {
	return "", "10", nil
}

func NewCsvQuizStore() *CsvQuizStore {
	return &CsvQuizStore{}
}
