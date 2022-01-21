package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvQuizStore struct {
	reader *csv.Reader
}

func (s *CsvQuizStore) NextQuiz() Quiz {
	next := Quiz{}
	record, err := s.reader.Read()
	if err != nil {
		next.Err = err
	} else if len(record) != 2 {
		next.Err = fmt.Errorf("wrong number of fields in a quiz record: %s. Expected 2 but got %d\n.", record, len(record))
	} else {
		next.Question, next.Answer = record[0], record[1]
	}
	return next
}

// Factory function which takes a describtor of csv filename
// and returns a pointer to CsvQuizStore
func NewCsvQuizStore(fd *os.File) *CsvQuizStore {
	return &CsvQuizStore{reader: csv.NewReader(fd)}
}
