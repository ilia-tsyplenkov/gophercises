package urlshort_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ilia-tsyplenkov/gophercises/urlshort"
	"github.com/ilia-tsyplenkov/gophercises/urlshort/test_sugar"
	"gopkg.in/yaml.v2"
)

var fallbackHandler = http.HandlerFunc(testFallbackHandler)
var doneCh = make(chan struct{})

const (
	dbFile              = "test.db"
	bucket              = "redirects"
	interval            = 10 * time.Millisecond
	fallbackHandlerBody = "fallback handler called"
)

func TestHandlerRedirectRequests(t *testing.T) {

	testCases := map[string]string{
		"/google":   "https://google.com",
		"/yandex":   "https://yandex.com",
		"/youtube":  "https://youtube.com",
		"/net/http": "https://pkg.go.dev/net/http",
	}

	handler := urlshort.MapHandler(testCases, nil)
	for from, to := range testCases {
		performRedirect(t, handler, from, to)

	}

}

func TestFallbackCalledNonRedirectRequests(t *testing.T) {
	handler := urlshort.MapHandler(nil, fallbackHandler)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected to have 200 code but got %d\n", w.Code)
	}
	got, err := w.Body.ReadString('\n')

	want := fallbackHandlerBody
	if err != nil {
		t.Fatalf("unxpected error while reading resonse body: %s\n", err)
	}
	got = strings.TrimSuffix(got, "\n")
	if got != want {
		t.Errorf("got %q in response body, but %q expected\n", got, want)
	}
}

func TestYAMLHandlerRedirectRequests(t *testing.T) {

	testCases := []struct {
		From string
		To   string
	}{
		{"/google", "https://google.com"},
		{"/yandex", "https://yandex.com"},
		{"/youtube", "https://youtube.com"},
		{"/net/http", "https://pkg.go.dev/net/http"},
	}
	yamlBinary, err := yaml.Marshal(testCases)
	if err != nil {
		t.Fatalf("error marshaling: %s\n", err)
	}

	handler, err := urlshort.YAMLHandler(yamlBinary, nil)
	if err != nil {
		t.Fatalf("error creating handler: %s\n", err)
	}
	for _, tc := range testCases {
		performRedirect(t, handler, tc.From, tc.To)

	}

}

func TestJsonHandlerRedirectRequests(t *testing.T) {

	testCases := []struct {
		From string
		To   string
	}{
		{"/google", "https://google.com"},
		{"/yandex", "https://yandex.com"},
		{"/youtube", "https://youtube.com"},
		{"/net/http", "https://pkg.go.dev/net/http"},
	}
	jsonBinary, err := json.Marshal(testCases)
	if err != nil {
		t.Fatalf("error marshaling: %s\n", err)
	}

	handler, err := urlshort.JSONHandler(jsonBinary, nil)
	if err != nil {
		t.Fatalf("error creating handler: %s\n", err)
	}
	for _, tc := range testCases {
		performRedirect(t, handler, tc.From, tc.To)

	}

}

func TestBoltDbHandlerRedirectRequests(t *testing.T) {

	testCases := map[string]string{
		"/google":   "https://google.com",
		"/yandex":   "https://yandex.com",
		"/youtube":  "https://youtube.com",
		"/net/http": "https://pkg.go.dev/net/http",
	}

	err := test_sugar.FillBucket(dbFile, bucket, testCases)
	if err != nil {
		t.Fatalf("error preparing test bucket - %q", err)
	}
	defer os.Remove(dbFile)
	handler, err := urlshort.BoltDbHandler(dbFile, fallbackHandler, interval, doneCh)
	if err != nil {
		t.Fatalf("error creating handler: %s\n", err)
	}
	defer func() {
		doneCh <- struct{}{}
	}()
	for from, to := range testCases {
		performRedirect(t, handler, from, to)

	}

}

func TestBoltDbHandlerUpToDateRedirects(t *testing.T) {

	testCases := map[string]string{
		"/google":   "https://google.com",
		"/yandex":   "https://yandex.com",
		"/youtube":  "https://youtube.com",
		"/net/http": "https://pkg.go.dev/net/http",
	}

	err := test_sugar.FillBucket(dbFile, bucket, testCases)
	if err != nil {
		t.Fatalf("error preparing test bucket - %q", err)
	}
	defer os.Remove(dbFile)
	handler, err := urlshort.BoltDbHandler(dbFile, fallbackHandler, interval, doneCh)
	if err != nil {
		t.Fatalf("error creating handler: %s\n", err)
	}
	defer func() {
		doneCh <- struct{}{}
	}()

	testCases["/go"] = "https://go.dev"
	err = test_sugar.FillBucket(dbFile, bucket, testCases)
	if err != nil {
		t.Fatalf("error updating test bucket - %q", err)
	}
	time.Sleep(10 * interval)
	for from, to := range testCases {
		performRedirect(t, handler, from, to)

	}

}
func testFallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, fallbackHandlerBody)
}

func performRedirect(t *testing.T, handler http.Handler, from, to string) {
	testName := fmt.Sprintf("RedirectFrom_%s", from)
	t.Run(testName, func(t *testing.T) {
		r := httptest.NewRequest("GET", from, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusFound {
			t.Fatalf("expected %d code but got %d\n", http.StatusFound, w.Code)
		}
		location := w.HeaderMap["Location"]
		if len(location) == 0 {
			t.Fatalf("redirecting from %q to %q: no 'Location' header empty or missed in the redirect response\n", from, to)
		}
		got := location[0]
		want := to

		if got != want {
			t.Fatalf("redirect from %q expected have %q location after redirect, but got %q\n", from, want, got)
		}

	})
}
