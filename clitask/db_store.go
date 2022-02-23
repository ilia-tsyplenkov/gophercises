package clitask

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	bolt "go.etcd.io/bbolt"
)

type DbStore struct {
	DbFile string
	Bucket string
}

func NewDbStore(dbFile, bucket string) *DbStore {
	store := &DbStore{DbFile: dbFile, Bucket: bucket}
	store.createBucket()
	return store
}

func (s *DbStore) ToDo() ([]Task, error) {
	db, err := bolt.Open(s.DbFile, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tasks := make([]Task, 0)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.Bucket))
		if b == nil {
			return fmt.Errorf("%q bucket wasn't found.", s.Bucket)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			task := Task{}
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			if !task.Done {
				tasks = append(tasks, task)
			}

		}
		return nil

	})
	if err != nil {
		return tasks, err
	}
	return tasks, nil

}

func (s *DbStore) Add(taskMsg string) error {

	db, err := bolt.Open(s.DbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	task := Task{Name: taskMsg}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.Bucket))
		if err != nil {
			return err
		}
		id, _ := b.NextSequence()
		task.Id = int(id)
		buf, err := json.Marshal(task)
		err = b.Put(Itob(task.Id), buf)
		return err
	})
	return err
}

func (s *DbStore) Update(task Task) error {
	db, err := bolt.Open(s.DbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.Bucket))
		if b == nil {
			return fmt.Errorf("%s bucket wasn't found.", s.Bucket)
		}
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}
		err = b.Put(Itob(task.Id), buf)
		return err
	})
	return err

}

func (s *DbStore) Do(id int) (Task, error) {
	tasks, err := s.ToDo()
	if err != nil {
		return Task{}, err
	}
	if id < 1 {
		return Task{}, ErrUnexpectedId
	}
	if id > len(tasks) {
		return Task{}, io.ErrUnexpectedEOF
	}
	id--
	todo := tasks[id]
	todo.Done = true
	err = s.Update(todo)
	return todo, err

}

func (s *DbStore) createBucket() error {
	db, err := bolt.Open(s.DbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.Bucket))

		return err
	})
	return err
}

// itob returns an 8-byte big endian representation of v.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
