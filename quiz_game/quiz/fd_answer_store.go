package quiz

import (
	"bufio"
	"os"
	"strings"
)

type FileAnswerStore struct {
	buffer *bufio.Reader
}

func NewFileAnswerStore(fileName string) (*FileAnswerStore, error) {
	var fd *os.File
	var err error
	if fileName == "" {
		fd = os.Stdin
	} else {
		fd, err = os.Open(fileName)
		if err != nil {
			return nil, err
		}

	}

	return &FileAnswerStore{buffer: bufio.NewReader(fd)}, nil
}

func (s *FileAnswerStore) NextAnswer() (string, error) {
	answer, err := s.buffer.ReadString('\n')
	if err != nil {
		return "", err
	}
	answer = strings.TrimSuffix(answer, "\n")

	return answer, nil
}
