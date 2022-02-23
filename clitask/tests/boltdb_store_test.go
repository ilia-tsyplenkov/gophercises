package external_tests

import (
	"encoding/json"
	"os"
	"testing"

	ct "github.com/gophercises/clitask"
	bolt "go.etcd.io/bbolt"
)

const (
	testDbFile = "testDb.db"
	testBucket = "tasks"
)

func TestDbToDoList(t *testing.T) {

	testCases := []struct {
		name     string
		dbfile   string
		bucket   string
		data     []ct.Task
		hasError bool
	}{
		{name: "wrongBucket", dbfile: testDbFile, bucket: "foo", data: []ct.Task{}, hasError: true},
		{name: "noTasks", dbfile: testDbFile, bucket: testBucket, data: []ct.Task{}, hasError: false},
		{name: "oneTask", dbfile: testDbFile, bucket: testBucket, data: []ct.Task{ct.Task{Name: "write test"}}, hasError: false},
		{name: "twoTasks", dbfile: testDbFile, bucket: testBucket, data: []ct.Task{ct.Task{Name: "write test"}, ct.Task{Name: "write code"}}, hasError: false},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			err := prepareDb(testDbFile, testBucket, tc.data)
			defer os.Remove(testDbFile)
			if err != nil {
				t.Fatalf("error preparing test database: %s\n", err)
			}
			store := ct.DbStore{DbFile: tc.dbfile, Bucket: tc.bucket}
			got, err := store.ToDo()
			if tc.hasError && err == nil {
				t.Fatalf("expected to have non-nil error but got %v\n", err)
			} else if !tc.hasError && err != nil {
				t.Fatalf("got unexpected error %v\n", err)
			}
			if !tasksEqual(got, tc.data) {
				t.Fatalf("expected to have %v task list, but got %v\n", tc.data, got)
			}
		})
	}
}

func TestDbStoreAdd(t *testing.T) {
	testCases := []struct {
		name         string
		newTask      string
		expectedInDb []ct.Task
	}{
		{name: "AddToEmptyDb", newTask: "write test", expectedInDb: []ct.Task{ct.Task{Id: 0, Name: "write test"}}},
		{name: "AddToExistingDbi", newTask: "write code", expectedInDb: []ct.Task{ct.Task{Id: 0, Name: "write test"}, ct.Task{Id: 1, Name: "write code"}}},
	}
	store := ct.DbStore{DbFile: testDbFile, Bucket: testBucket}
	defer os.Remove(testDbFile)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store.Add(tc.newTask)
			got, err := store.ToDo()
			if err != nil {
				t.Fatalf("unexpected error - %v\n", err)
			}
			if !tasksEqual(got, tc.expectedInDb) {
				t.Fatalf("expected to have %v task list, but got %v\n", tc.expectedInDb, got)
			}
		})

	}
}

func TestDbStoreDo(t *testing.T) {
	tasksInDb := getTasksFromNames("write test", "write code", "pass test")
	err := prepareDb(testDbFile, testBucket, tasksInDb)
	defer os.Remove(testDbFile)
	if err != nil {
		t.Fatalf("error preparing test database: %s\n", err)
	}
	store := ct.DbStore{DbFile: testDbFile, Bucket: testBucket}
	store.Do(1)
	got, err := store.ToDo()
	if err != nil {
		t.Fatalf("error getting tasks from db: %s\n", err)
	}
	want := getTasksFromNames("write code", "pass test")
	if !tasksEqual(got, want) {
		t.Fatalf("expected to have next task list:\n%v\nbut got:\n%v\n", want, got)
	}

}

func TestDbStoreUpdateTask(t *testing.T) {
	task := ct.Task{Id: 1, Name: "write test", Done: false}
	err := prepareDb(testDbFile, testBucket, []ct.Task{task})
	defer os.Remove(testDbFile)
	if err != nil {
		t.Fatalf("error preparing test database: %s\n", err)
	}
	store := ct.DbStore{DbFile: testDbFile, Bucket: testBucket}
	task.Done = true
	store.Update(task)
	tasks, err := store.ToDo()
	if err != nil {
		t.Fatalf("error getting tasks from db: %s\n", err)
	}
	if len(tasks) > 0 {
		t.Fatalf("expected to have 0 tasks to complete, but got: %v\n", tasks)
	}
}

func getTasksFromNames(names ...string) []ct.Task {
	tasks := make([]ct.Task, 0)
	for id, name := range names {
		task := ct.Task{Id: id + 1, Name: name}
		tasks = append(tasks, task)
	}
	return tasks

}
func tasksEqual(got, want []ct.Task) bool {
	if len(got) != len(want) {
		return false
	}
	for i, v := range got {
		if v.Name != want[i].Name || v.Done || want[i].Done {
			return false
		}
	}
	return true
}

func prepareDb(dbFile string, bucket string, data []ct.Task) error {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return err
	}

	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		for _, task := range data {
			id, _ := b.NextSequence()
			task.Id = int(id)
			buf, err := json.Marshal(task)
			err = b.Put(ct.Itob(task.Id), buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err

}
