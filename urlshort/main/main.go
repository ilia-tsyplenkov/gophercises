package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ilia-tsyplenkov/gophercises/urlshort"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Page %q handled by default handler\n", r.URL.Path)
}

var roadFile string

func init() {
	flag.StringVar(&roadFile, "road_src", "redirects.yaml", "source file with redirect rules")
}
func main() {
	flag.Parse()
	fd, err := os.Open(roadFile)
	if err != nil {
		log.Fatalf("cannot open road source file: %s\n", err)
	}
	defer fd.Close()
	yamlBytes, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Fatalf("error while reading roadmap data: %s\n", err)
	}

	redirectHandler, err := urlshort.YAMLHandler(yamlBytes, http.HandlerFunc(defaultHandler))
	http.Handle("/", redirectHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
