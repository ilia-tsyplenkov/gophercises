package clitask

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestReadCommand(t *testing.T) {
	file := "testInput"
	testCases := []struct {
		name string
		cmd  string
		want string
		err  error
	}{
		{
			name: "EmptyCmd", cmd: "", want: "", err: io.EOF,
		},

		{
			name: "TaskCmd", cmd: "task\n", want: "task", err: nil,
		},
		{
			name: "TaskListCmd", cmd: "task list\n", want: "task list", err: nil,
		},
		{
			name: "TaskDoCmd", cmd: "task do 1\n", want: "task do 1", err: nil,
		},
		{
			name: "TaskAddCmd", cmd: "task add foo\n", want: "task add foo", err: nil,
		},
		{
			name: "UnknownFooCmd", cmd: "task foo\n", want: "", err: ErrUnknownCmd,
		},
		{
			name: "UnknownSpamBarCmd", cmd: "task spam bar\n", want: "", err: ErrUnknownCmd,
		},
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
			if got != tc.want {
				t.Fatalf("expected to have %q command from input, but got - %q\n", tc.want, got)
			}
		})
	}
}

func TestFixCommand(t *testing.T) {
	testCases := []struct {
		name string
		cmd  string
		want string
	}{
		{name: "task", cmd: "task", want: "task"},
		{name: "taskWithNLCharacter", cmd: "task\n", want: "task"},
		{name: "taskWithSomeSpaces", cmd: "  task  ", want: "task"},
		{name: "taskDo", cmd: "task do 1", want: "task do 1"},
		{name: "taskDoWithSpacesBetweenWords", cmd: "task   do  1", want: "task do 1"},
	}
	manager := Manager{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := manager.fixCmd(tc.cmd)
			if got != tc.want {
				t.Fatalf("expected to have %q after fix, but got %q\n", tc.want, got)
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
		{name: "taskDo", cmd: "task do 1", known: true},
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
