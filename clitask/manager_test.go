package clitask

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

var ErrUnknownCmd = errors.New("unsupported command.")

func TestReadCommand(t *testing.T) {
	file := "testInput"
	testCases := []struct {
		name string
		cmd  string
		err  error
	}{
		{
			name: "EmptyCmd", cmd: "", err: io.EOF,
		},

		{
			name: "TaskCmd", cmd: "task\n", err: nil,
		},
		{
			name: "TaskListCmd", cmd: "task list\n", err: nil,
		},
		{
			name: "TaskDoCmd", cmd: "task do 1\n", err: nil,
		},
		{
			name: "TaskAddCmd", cmd: "task add foo\n", err: nil,
		},
		{
			name: "TaskAddCmd", cmd: "task add foo\n", err: nil,
		},
		// {
		// 	name: "UnknownFooCmd", cmd: "task foo\n", err: ErrUnknownCmd,
		// },
		// {
		// 	name: "UnknownSpamBarCmd", cmd: "task spam bar\n", err: ErrUnknownCmd,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fd, err := os.Create(file)
			if err != nil {
				t.Fatalf("error create file: %s\n", err)
			}
			defer func() {
				fd.Close()
				os.Remove(file)
			}()
			fmt.Fprint(fd, tc.cmd)
			fd.Seek(0, 0)
			manager := Manager{fd}
			got, err := manager.ReadCmd()
			if err != tc.err {
				t.Fatalf("expected to have %q error but got %q", tc.err, err)
			}
			if got != tc.cmd {
				t.Fatalf("expected to have %q command from input, but got - %q\n", tc.cmd, got)
			}
		})
	}
}

func TestIsKnownCommand(t *testing.T) {
	testCases := []struct {
		name  string
		cmd   string
		known bool
	}{
		{name: "task", cmd: "task", known: true},
		{name: "taskWithSomeSpaces", cmd: "  task  ", known: true},
		{name: "taskDo", cmd: "task do 1", known: true},
		{name: "taskDoWithSpacesBetweenWords", cmd: "task   do  1", known: true},
		{name: "taskAdd", cmd: "task add foo", known: true},
		{name: "taskList", cmd: "task list", known: true},
		{name: "taskFoo", cmd: "task foo", known: false},
		{name: "taskSpamBar", cmd: "task spam bar", known: false},
	}
	manager := Manager{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := manager.isKnown(tc.cmd)
			if got != tc.known {
				t.Fatalf("expected to have %v known status of command but got: %v", tc.known, got)
			}

		})
	}
}
