package external_tests

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/clitask"
)

func TestWork(t *testing.T) {

	testCases := []struct {
		name    string
		command string
		tasks   []string
		want    string
		err     error
	}{
		{name: "help", command: "task\n", want: clitask.Greeting + clitask.HelpMsg},
		{name: "listEmptyTask", command: "task list\n", tasks: nil, want: clitask.Greeting + clitask.EmptyBacklog},
		{name: "listOneTask", command: "task list\n", tasks: []string{"write test"}, want: clitask.Greeting + "1. write test\n"},
		{name: "listTwoTasks", command: "task list\n", tasks: []string{"write test", "write code"}, want: clitask.Greeting + "1. write test\n2. write code\n"},
		{name: "addOneTask", command: "task add write code\n", want: clitask.Greeting + "Added \"write code\" to your task list.\n"},
		{name: "doOneTask", command: "task do 1\n", tasks: []string{"write test", "write code"}, want: clitask.Greeting + "You have completed the \"write test\" task.\n"},
		{name: "unknownCommand", command: "task foo\n", want: clitask.Greeting, err: clitask.ErrUnknownCmd},
		{name: "doTaskUnexistingId", command: "task do 10\n", want: clitask.Greeting, err: io.ErrUnexpectedEOF},
		{name: "doTaskZeroId", command: "task do 0\n", want: clitask.Greeting, err: clitask.ErrUnexpectedId},
		{name: "doTaskUnapplicableId", command: "task do foo\n", want: clitask.Greeting, err: clitask.ErrIncorrectId},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				os.Remove("testInput")
				os.Remove("testOutput")
			}()
			in, err := os.Create("testInput")
			if err != nil {
				t.Fatal("error creating test input file")
			}
			defer in.Close()
			out, err := os.Create("testOutput")
			if err != nil {
				t.Fatal("error creating test output file")
			}
			defer out.Close()
			in.WriteString(tc.command)
			in.Seek(0, 0)
			store := clitask.NewMemStore()
			for _, task := range tc.tasks {
				store.Add(task)
			}
			manager := clitask.Manager{Input: in, Output: out, Store: store}
			err = manager.Work()
			if err != tc.err {
				t.Fatalf("expected to have next error - %v, but got - %v", tc.err, err)
			}

			out.Seek(0, 0)
			res, err := ioutil.ReadAll(out)
			if err != nil {
				t.Fatalf("error reading test output: %s\n", err)
			}
			got := string(res)
			if got != tc.want {
				t.Fatalf("expected to have %q as worker result, but got %q\n", tc.want, got)
			}
		})
	}
}
