package clitask

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetTask(t *testing.T) {
	want := "write tests"
	got := GetTask()
	if got != want {
		t.Fatalf("expected to have %q task but got %q\n", want, got)
	}
}

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

	have := []Task{Task{name: "write test", done: false}}
	store := MemStore{have, nil}
	store.Do(1)
	want := []Task{Task{name: "write test", done: true}}
	got := store.Completed()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected to have %v completed tasks but got %v", want, got)
	}

}
