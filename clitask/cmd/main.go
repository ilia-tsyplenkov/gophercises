package main

import (
	"io"
	"log"
	"os"

	"github.com/ilia-tsyplenkov/gophercises/clitask"
)

func main() {
	dbFile := "tasks.db"
	tasksBucket := "tasks"

	store := clitask.NewDbStore(dbFile, tasksBucket)
	manager := clitask.NewManager(os.Stdin, os.Stdout, store)

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
