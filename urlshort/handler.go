package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

type RedirectHandler struct {
	redirects map[string]string
	fallback  http.Handler
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	value, ok := h.redirects[r.URL.Path]
	if ok {
		http.Redirect(w, r, value, http.StatusFound)
		return
	}
	h.fallback.ServeHTTP(w, r)
}

func MapHandler(redirects map[string]string, fallback http.Handler) http.Handler {
	return &RedirectHandler{redirects, fallback}
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.Handler, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.Handler, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func BoltDbHandler(dbFile string, fallback http.Handler) (http.Handler, error) {
	buildMap, err := readDb(dbFile, "redirects")
	if err != nil {
		return nil, err
	}
	handler := MapHandler(buildMap, fallback)
	hd, ok := handler.(*RedirectHandler)
	if !ok {
		return handler, nil
	}
	go func() {
		time.Sleep(1 * time.Second)
		redirects, _ := readDb(dbFile, "redirects")
		log.Println("fetched redirects from goroutine:", redirects)
		hd.redirects = redirects
	}()
	return hd, nil
}

type redirect struct {
	From string
	To   string
}

func buildMap(d []redirect) map[string]string {
	res := make(map[string]string)

	for _, item := range d {
		res[item.From] = item.To
	}
	return res

}

func parseYAML(yamlBinary []byte) ([]redirect, error) {
	res := make([]redirect, 0)
	err := yaml.Unmarshal(yamlBinary, &res)
	return res, err
}

func parseJSON(binary []byte) ([]redirect, error) {
	res := make([]redirect, 0)
	err := json.Unmarshal(binary, &res)
	return res, err
}

func readDb(dbFile, bucketName string) (map[string]string, error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	res := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("%s bucket wasn't found.", bucketName)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			res[string(k)] = string(v)

		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
