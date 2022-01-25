package short_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	url "github.com/ilia-tsyplenkov/gophercises/urlshort"
)

func TestRedirectHandlerLocation(t *testing.T) {

	testCases := map[string]string{
		"/google":   "https://google.com",
		"/yandex":   "https://yandex.com",
		"/youtube":  "https://youtube.com",
		"/net/http": "https://pkg.go.dev/net/http",
	}

	handler := url.RedirectHandler{testCases}
	for from, to := range testCases {

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
	}

}
