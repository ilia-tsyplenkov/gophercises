package clitask

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var ErrUnknownCmd = errors.New("unsupported command.")
var knownCommand = []string{"task", "task list", "task do", "task add"}

type Manager struct {
	input io.Reader
}

func (m *Manager) ReadCmd() (string, error) {
	reader := bufio.NewReader(m.input)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !m.isKnown(s) {
		return "", ErrUnknownCmd
	}
	return s, nil

}

func (m *Manager) isKnown(s string) bool {
	s = strings.Trim(s, " \n")
	fixed := ""
	// remove all unnecessary spaces between words
	for _, ch := range s {
		if ch == ' ' && fixed != "" && fixed[len(fixed)-1] == ' ' {
			continue
		}
		fixed = fixed + string(ch)
	}
	var cmd string

	parts := strings.Split(fixed, " ")
	if len(parts) == 1 {
		cmd = parts[0]
	} else {
		cmd = strings.Join(parts[:2], " ")
	}

	for _, item := range knownCommand {
		if cmd == item {
			return true
		}
	}
	return false
}
