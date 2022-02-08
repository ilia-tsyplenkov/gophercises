package main

import (
	"io"
	"log"
	"os"

	"github.com/gophercises/clitask"
)

func main() {
	manager := clitask.NewManager(os.Stdin, os.Stdout)
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
