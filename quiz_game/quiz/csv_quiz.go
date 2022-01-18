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

func NewCsvQuizStore(filePath string) (*CsvQuizStore, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// defer fd.Close()
	store := CsvQuizStore{reader: csv.NewReader(fd)}
	return &store, nil
}
