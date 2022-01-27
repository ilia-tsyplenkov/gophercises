package urlshort

import (
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

func TestParseYAML(t *testing.T) {
	originData := []redirect{
		{"/python", "https://python.org"},
		{"/go", "https://golang.org"},
	}

	yamlData, err := yaml.Marshal(originData)
	if err != nil {
		t.Fatalf("error yaml marhaling: %s\n", err)
	}
	got, err := parseYAML(yamlData)
	want := originData
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("exptected to have %q after parsing but got %q\n", want, got)
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
