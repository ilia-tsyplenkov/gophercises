package clitask

import (
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

func TestGetTaskList(t *testing.T) {
	want := []string{"write test", "write code", "pass test"}
	store := MemStore{want}
	got := store.TaskList()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected to have %q tasks list but got %q\n", want, got)
	}

}
