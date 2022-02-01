package urlshort

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestBuildingMap(t *testing.T) {
	testFile := "buildMap.yaml"
	originData := []redirect{
		{"/python", "https://python.org"},
		{"/go", "https://golang.org"},
	}

	err := createYamlFile(testFile, originData)
	if err != nil {
		t.Fatalf(err.Error())
	}

	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		f.Close()
		os.Remove(testFile)
	}()
	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf(err.Error())
	}
	data := make([]redirect, 0)
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		t.Fatal(err)
	}
	got := buildMap(data)
	want := map[string]string{
		"/python": "https://python.org",
		"/go":     "https://golang.org",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected to have %q, but got %q\n", want, got)
	}
}

func TestParser(t *testing.T) {
	originData := []redirect{
		{"/python", "https://python.org"},
		{"/go", "https://golang.org"},
	}

	testCases := []struct {
		Name        string
		MarshalFunc func(interface{}) ([]byte, error)
		ParseFunc   func([]byte) ([]redirect, error)
	}{
		{"YAML", yaml.Marshal, parseYAML},
		{"JSON", json.Marshal, parseJSON},
	}

	for _, tc := range testCases {
		testName := "parse" + tc.Name
		t.Run(testName, func(t *testing.T) {
			data, err := tc.MarshalFunc(originData)
			if err != nil {
				t.Fatalf("error yaml marhaling: %s\n", err)
			}

			got, err := tc.ParseFunc(data)
			if err != nil {
				t.Fatalf("error parsing data: %s\n", err)
			}
			want := originData
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("exptected to have %q after parsing but got %q\n", want, got)
			}
		})
	}

}

func TestReadDb(t *testing.T) {
	pathMap, err := readDb("test.db", "redirects")
	if err != nil {
		t.Fatalf("expected to have success read from db, but got %q\n", err)
	}
	if len(pathMap) == 0 {
		t.Fatal("expected to read non 0 record from db.")
	}

}

func createYamlFile(fileName string, data interface{}) error {
	binaryData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, binaryData, 0644)
	if err != nil {
		return err
	}
	return nil
}
