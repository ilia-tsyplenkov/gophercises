package urlshort_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	url "github.com/ilia-tsyplenkov/gophercises/urlshort"
)

var fallbackHandlerBody string = "fallback handler called"

func TestHandlerRedirectRequests(t *testing.T) {

	testCases := map[string]string{
		"/google":   "https://google.com",
		"/yandex":   "https://yandex.com",
		"/youtube":  "https://youtube.com",
		"/net/http": "https://pkg.go.dev/net/http",
	}

	handler := url.MapHandler(testCases, nil)
	for from, to := range testCases {
		testName := fmt.Sprintf("RedirectFrom_%s", from)
		t.Run(testName, func(t *testing.T) {
			r := httptest.NewRequest("GET", from, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			if w.Code != http.StatusFound {
				t.Fatalf("expected %d code but got %d\n", w.Code, http.StatusFound)
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

}

func TestFallbackCalledNonRedirectRequests(t *testing.T) {
	fallbackHandler := http.HandlerFunc(testFallbackHandler)
	handler := url.MapHandler(nil, fallbackHandler)
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

func testFallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, fallbackHandlerBody)
}