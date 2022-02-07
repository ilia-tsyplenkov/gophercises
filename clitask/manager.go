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
var knownCommand = []string{"task", "task list", "task do", "task add"}

type Manager struct {
	input  io.Reader
	output io.Writer
	store  *MemStore
}

func (m *Manager) ReadCmd() (string, error) {
	reader := bufio.NewReader(m.input)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	s = m.fixCmd(s)
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

func (m *Manager) Work() {
	fullCmd, err := m.ReadCmd()
	if err != nil {
		panic(err)
	}
	cmd, args := m.splitOnArgs(fullCmd)
	switch cmd {
	case "task":
		m.WriteResult(helpMsg)
	case "task list":
		tasks := m.store.ToDo()
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
			panic(err)
		}
		task, err := m.store.Do(id)
		if err != nil {
			panic(err)
		}
		m.WriteResult(fmt.Sprintf("You have completed the %q task.\n", task.name))
	}

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
