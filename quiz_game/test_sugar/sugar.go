package test_sugar

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func MakeTestAnswerFile(fileName string, answers []string) error {
	fd, err := os.Create(fileName)
	if err == os.ErrExist {
		os.Remove(fileName)
		fd, err = os.Create(fileName)
		if err != nil {
			return err
		}
	}
	defer fd.Close()
	for _, a := range answers {
		fmt.Fprintf(fd, "%s\n", a)
	}
	return nil

}

func MakeTestCsv(filename string, records [][]string) (filePath string, answers []string) {

	fd, err := os.Create(filename)
	if err == os.ErrExist {
		os.Remove(filename)
		fd, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
	defer fd.Close()
	answers = make([]string, len(records))
	w := csv.NewWriter(fd)
	for i, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error preparing csv test file:", err)
		}
		answers[i] = record[1]
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return fd.Name(), answers

}
