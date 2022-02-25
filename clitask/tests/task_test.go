package external_tests

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	ct "github.com/ilia-tsyplenkov/gophercises/clitask"
)

func TestTaskToDoList(t *testing.T) {

	testCases := []struct {
		have []ct.Task
		want []ct.Task
	}{
		{
			have: nil, want: nil,
		},
		{
			have: []ct.Task{ct.Task{Name: "write test"}},
			want: []ct.Task{ct.Task{Name: "write test"}},
		},
		{
			have: []ct.Task{ct.Task{Name: "write test"}, ct.Task{Name: "write code"}},
			want: []ct.Task{ct.Task{Name: "write test"}, ct.Task{Name: "write code"}},
		},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%dtasks", len(tc.have))
		t.Run(testName, func(t *testing.T) {
			store := ct.MemStore{Todo: tc.have}
			got, _ := store.ToDo()
			if !reflect.DeepEqual(got, tc.have) {
				t.Fatalf("expected to have %v tasks list but got %v\n", tc.have, got)
			}
		})
	}

}

func TestAddNewTask(t *testing.T) {
	have := []ct.Task{
		ct.Task{Name: "write test", Done: false},
		ct.Task{Name: "write code", Done: false},
	}
	newTask := ct.Task{Name: "pass test", Done: false}

	store := ct.MemStore{have, nil}
	want := make([]ct.Task, len(have))
	copy(want, have)
	want = append(want, newTask)
	store.Add(newTask.Name)
	got, _ := store.ToDo()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected to have %v tasks list but got %v\n", want, got)
	}

}

func TestCompletedTaskList(t *testing.T) {
	testCases := []struct {
		have []ct.Task
		want []ct.Task
	}{
		{
			have: nil, want: nil,
		},
		{
			have: []ct.Task{ct.Task{Name: "write test", Done: true}},
			want: []ct.Task{ct.Task{Name: "write test", Done: true}},
		},
		{
			have: []ct.Task{ct.Task{Name: "write test", Done: true}, ct.Task{Name: "write code", Done: true}},
			want: []ct.Task{ct.Task{Name: "write test", Done: true}, ct.Task{Name: "write code", Done: true}},
		},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%dtasks", len(tc.have))
		t.Run(testName, func(t *testing.T) {
			store := ct.MemStore{Done: tc.have}
			got := store.Completed()
			if !reflect.DeepEqual(got, tc.have) {
				t.Fatalf("expected to have %v tasks list but got %v\n", tc.have, got)
			}
		})
	}

}

func TestDoTask(t *testing.T) {
	testCases := []struct {
		name string
		have []ct.Task
		want []ct.Task
		toDo int
		err  error
	}{
		{
			name: "EmptyBacklog",
			have: nil,
			want: nil,
			toDo: 1,
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "1TaskInBacklogButId2Requested",
			have: []ct.Task{ct.Task{Name: "write test", Done: false}},
			want: nil,
			toDo: 2,
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "1TaskInBacklogButId0Requested",
			have: []ct.Task{ct.Task{Name: "write test", Done: false}},
			want: nil,
			toDo: 0,
			err:  ct.ErrUnexpectedId,
		},
		{
			name: "1TaskId1Requested",
			have: []ct.Task{ct.Task{Name: "write test", Done: false}},
			want: []ct.Task{ct.Task{Name: "write test", Done: true}},
			toDo: 1,
			err:  nil,
		},
		{
			name: "2TasksId1Requested",
			have: []ct.Task{ct.Task{Name: "write test", Done: false}, ct.Task{Name: "write code", Done: false}},
			want: []ct.Task{ct.Task{Name: "write test", Done: true}},
			toDo: 1,
			err:  nil,
		},
		{
			name: "2TasksId2Requested",
			have: []ct.Task{ct.Task{Name: "write test", Done: false}, ct.Task{Name: "write code", Done: false}},
			want: []ct.Task{ct.Task{Name: "write code", Done: true}},
			toDo: 2,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			todoNum := len(tc.have)
			store := ct.MemStore{Todo: tc.have}
			_, err := store.Do(tc.toDo)
			if err != tc.err {
				t.Fatalf("expected to have %v error but got %v\n", tc.err, err)
			}
			got := store.Completed()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected to have %v in completed but got %v\n", tc.want, got)
			}
			if tc.want != nil {
				t.Log("Check that ToDo backlog has been reduced.")
				tasks, _ := store.ToDo()
				gotTodoNum := len(tasks)
				wantTodoNum := todoNum - 1
				if gotTodoNum != wantTodoNum {
					t.Fatalf("expected to have %d tasks in backlog but got %d\n", wantTodoNum, gotTodoNum)
				}
			}

		})
	}
}
