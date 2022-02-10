package clitask

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type DbStore struct {
	DbFile string
	Bucket string
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
			return fmt.Errorf("%s bucket wasn't found.", s.Bucket)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			task := Task{}
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			tasks = append(tasks, task)

		}
		return nil

	})
	if err != nil {
		return tasks, err
	}
	return tasks, nil

}
