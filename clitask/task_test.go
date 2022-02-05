package clitask

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestTaskToDoList(t *testing.T) {

	testCases := []struct {
		have []Task
		want []Task
	}{
		{
			have: nil, want: nil,
		},
		{
			have: []Task{Task{name: "write test"}},
			want: []Task{Task{name: "write test"}},
		},
		{
			have: []Task{Task{name: "write test"}, Task{name: "write code"}},
			want: []Task{Task{name: "write test"}, Task{name: "write code"}},
		},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%dtasks", len(tc.have))
		t.Run(testName, func(t *testing.T) {
			store := MemStore{todo: tc.have}
			got := store.ToDo()
			if !reflect.DeepEqual(got, tc.have) {
				t.Fatalf("expected to have %v tasks list but got %v\n", tc.have, got)
			}
		})
	}

}

func TestAddNewTask(t *testing.T) {
	have := []Task{
		Task{name: "write test", done: false},
		Task{name: "write code", done: false},
	}
	newTask := Task{name: "pass test", done: false}

	store := MemStore{have, nil}
	want := make([]Task, len(have))
	copy(want, have)
	want = append(want, newTask)
	store.Add(newTask.name)
	got := store.ToDo()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected to have %v tasks list but got %v\n", want, got)
	}

}

func TestCompletedTaskList(t *testing.T) {
	testCases := []struct {
		have []Task
		want []Task
	}{
		{
			have: nil, want: nil,
		},
		{
			have: []Task{Task{name: "write test", done: true}},
			want: []Task{Task{name: "write test", done: true}},
		},
		{
			have: []Task{Task{name: "write test", done: true}, Task{name: "write code", done: true}},
			want: []Task{Task{name: "write test", done: true}, Task{name: "write code", done: true}},
		},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%dtasks", len(tc.have))
		t.Run(testName, func(t *testing.T) {
			store := MemStore{done: tc.have}
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
		have []Task
		want []Task
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
			have: []Task{Task{name: "write test", done: false}},
			want: nil,
			toDo: 2,
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "1TaskInBacklogButId0Requested",
			have: []Task{Task{name: "write test", done: false}},
			want: nil,
			toDo: 0,
			err:  ErrUnexpectedId,
		},
		{
			name: "1TaskId1Requested",
			have: []Task{Task{name: "write test", done: false}},
			want: []Task{Task{name: "write test", done: true}},
			toDo: 1,
			err:  nil,
		},
		{
			name: "2TasksId1Requested",
			have: []Task{Task{name: "write test", done: false}, Task{name: "write code", done: false}},
			want: []Task{Task{name: "write test", done: true}},
			toDo: 1,
			err:  nil,
		},
		{
			name: "2TasksId2Requested",
			have: []Task{Task{name: "write test", done: false}, Task{name: "write code", done: false}},
			want: []Task{Task{name: "write code", done: true}},
			toDo: 2,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			todoNum := len(tc.have)
			store := MemStore{todo: tc.have}
			err := store.Do(tc.toDo)
			if err != tc.err {
				t.Fatalf("expected to have %v error but got %v\n", tc.err, err)
			}
			got := store.Completed()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected to have %v in completed but got %v\n", tc.want, got)
			}
			if tc.want != nil {
				t.Log("Check that ToDo backlog has been reduced.")
				gotTodoNum := len(store.ToDo())
				wantTodoNum := todoNum - 1
				if gotTodoNum != wantTodoNum {
					t.Fatalf("expected to have %d tasks in backlog but got %d\n", wantTodoNum, gotTodoNum)
				}
			}

		})
	}
}
