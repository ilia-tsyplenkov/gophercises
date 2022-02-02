package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ilia-tsyplenkov/gophercises/urlshort"
)

var roadFile string
var defaultHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Page %q handled by default handler\n", r.URL.Path)
})

const (
	yamlType      = "yaml"
	jsonType      = "json"
	fetchInterval = 30 * time.Second
)

func init() {
	flag.StringVar(&roadFile, "src", "redirects.yaml", "source file with redirect rules")
}
func main() {
	var done chan struct{}
	flag.Parse()
	s := strings.Split(roadFile, ".")
	fileType := s[len(s)-1]
	log.Printf("provided source file: %s\n", roadFile)
	log.Printf("source encoding: %s\n", fileType)

	var handler http.Handler
	var binaryData []byte
	var err error
	if fileType == yamlType || fileType == jsonType {
		fd, err := os.Open(roadFile)
		if err != nil {
			log.Fatalf("cannot open road source file: %s\n", err)
		}
		defer fd.Close()
		binaryData, err = ioutil.ReadAll(fd)
		if err != nil {
			log.Fatalf("error while reading roadmap data: %s\n", err)
		}
	}

	switch fileType {
	case yamlType:
		log.Println("creating yaml handler")
		handler, err = urlshort.YAMLHandler(binaryData, defaultHandler)
	case jsonType:
		log.Println("creating json handler")
		handler, err = urlshort.JSONHandler(binaryData, defaultHandler)
	default:
		done = make(chan struct{})
		handler, err = urlshort.BoltDbHandler(roadFile, defaultHandler, fetchInterval, done)
		defer func() {
			done <- struct{}{}
		}()
	}
	if err != nil {
		log.Fatalf("error creating handler: %s\n", err)
	}
	http.Handle("/", handler)
	log.Println("Server is starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
