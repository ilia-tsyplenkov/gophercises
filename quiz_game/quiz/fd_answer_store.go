package quiz

import (
	"bufio"
	"os"
	"strings"
)

type FileAnswerStore struct {
	buffer *bufio.Reader
}

// Factory function which takes descriptor of file with answers
// and returns pointer to FileAnswerStore.
func NewFileAnswerStore(f *os.File) *FileAnswerStore {
	return &FileAnswerStore{buffer: bufio.NewReader(f)}
}

func (s *FileAnswerStore) NextAnswer() (string, error) {
	answer, err := s.buffer.ReadString('\n')
	if err != nil {
		return "", err
	}
	answer = strings.TrimSuffix(answer, "\n")

	return answer, nil
}
