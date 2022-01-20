package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvQuizStore struct {
	reader *csv.Reader
}

func (s *CsvQuizStore) NextQuiz() (string, string, error) {
	record, err := s.reader.Read()
	if err != nil {
		return "", "", err
	}
	if len(record) != 2 {
		return "", "", fmt.Errorf("wrong number of fields in a quiz record: %s. Expected 2 but got %d\n.", record, len(record))
	}
	return record[0], record[1], nil
}

// Factory function which takes a describtor of csv filename
// and returns a pointer to CsvQuizStore
func NewCsvQuizStore(fd *os.File) *CsvQuizStore {
	return &CsvQuizStore{reader: csv.NewReader(fd)}
}
