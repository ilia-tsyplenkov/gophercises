package main

import (
	"io"
	"log"
	"os"

	"github.com/gophercises/clitask"
)

func main() {
	// manager := clitask.NewManager(os.Stdin, os.Stdout)
	dbFile := "tasks.db"
	tasksBucket := "tasks"

	store := &clitask.DbStore{DbFile: dbFile, Bucket: tasksBucket}
	manager := &clitask.Manager{Input: os.Stdin, Output: os.Stdout, Store: store}
	for {
		err := manager.Work()
		if err != nil {
			log.Println(err)
			if err == io.EOF {
				break
			}
		}
	}
	log.Println("Program exiting.")
}
