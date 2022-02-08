package clitask

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ErrUnknownCmd = errors.New("unsupported command.")
var ErrIncorrectId = errors.New("incorrect task id. id must be positibe integer value.")
var knownCommand = []string{"task", "task list", "task do", "task add"}

const (
	EmptyBacklog string = "your backlog is empty.\n"
	greeting     string = "$ "
)

type Manager struct {
	input  io.Reader
	output io.Writer
	store  *MemStore
}

func NewManager(in io.Reader, out io.Writer) *Manager {
	return &Manager{input: in, output: out, store: NewMemStore()}
}

func (m *Manager) ReadCmd() (string, error) {
	reader := bufio.NewReader(m.input)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	s = m.fixCmd(s)
	// '\n' was pressed
	if s == "" {
		return s, nil
	}
	if !m.isKnown(s) {
		return "", ErrUnknownCmd
	}
	return s, nil

}

func (m *Manager) WriteResult(data interface{}) error {
	var err error

	switch d := data.(type) {
	case []string:
		for i, s := range d {
			s = fmt.Sprintf("%d. %s\n", i+1, s)
			if _, err = m.output.Write([]byte(s)); err != nil {
				return err
			}

		}
	case error:
		_, err = m.output.Write([]byte(d.Error()))
	case string:
		_, err = m.output.Write([]byte(d))
	default:
		err = fmt.Errorf("unexpected data type.")
	}

	return err
}

func (m *Manager) Work() error {
	m.output.Write([]byte(greeting))
	fullCmd, err := m.ReadCmd()
	if err != nil {
		return err
	}
	cmd, args := m.splitOnArgs(fullCmd)
	switch cmd {
	case "task":
		m.WriteResult(helpMsg)
	case "task list":
		tasks := m.store.ToDo()
		if len(tasks) == 0 {
			m.WriteResult(EmptyBacklog)
		}
		t := make([]string, 0)
		for _, task := range tasks {
			t = append(t, task.name)
		}
		m.WriteResult(t)
	case "task add":
		m.store.Add(args)
		m.WriteResult(fmt.Sprintf("Added %q to your task list.\n", args))
	case "task do":
		id, err := strconv.Atoi(args)
		if err != nil {
			return ErrIncorrectId
		}
		task, err := m.store.Do(id)
		if err != nil {
			return err
		}
		m.WriteResult(fmt.Sprintf("You have completed the %q task.\n", task.name))
	case "":
		{
		}
	default:
		return fmt.Errorf("command %q is known but not supported by worker yet.\n", cmd)
	}
	return nil

}

func (m *Manager) fixCmd(s string) string {
	s = strings.Trim(s, " \n")
	fixed := ""
	// remove all unnecessary spaces between words
	for _, ch := range s {
		if ch == ' ' && fixed != "" && fixed[len(fixed)-1] == ' ' {
			continue
		}
		fixed = fixed + string(ch)
	}
	return fixed
}

func (m *Manager) isKnown(s string) bool {
	var cmd string

	parts := strings.Split(s, " ")
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

func (m *Manager) splitOnArgs(s string) (string, string) {
	var cmd, args string
	parts := strings.Split(s, " ")
	if len(parts) == 1 {
		cmd, args = parts[0], ""
	} else {
		cmd, args = strings.Join(parts[:2], " "), strings.Join(parts[2:], " ")
	}
	return cmd, args
}
