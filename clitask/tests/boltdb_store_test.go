package external_tests

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"reflect"
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
		{name: "oneTask", dbfile: testDbFile, bucket: testBucket, data: []ct.Task{ct.Task{Id: 1, Name: "write test"}}, hasError: false},
		{name: "twoTasks", dbfile: testDbFile, bucket: testBucket, data: []ct.Task{ct.Task{Id: 1, Name: "write test"}, ct.Task{Id: 2, Name: "write code"}}, hasError: false},
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
			if !reflect.DeepEqual(got, tc.data) {
				t.Fatalf("expected to have %v task list, but got %v\n", tc.data, got)
			}
		})
	}
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
			err = b.Put(itob(task.Id), buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err

}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
