package url_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/short/url"
)

func TestRedirectHandlerStatusCode(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler := url.RedirectHandler{}
	handler.ServeHTTP(w, r)
	want := http.StatusFound
	got := w.Code
	if got != want {
		t.Fatalf("expected to have %d code but got %d\n", want, got)
	}
}

func TestRedirectHandlerLocation(t *testing.T) {

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler := url.RedirectHandler{}
	handler.ServeHTTP(w, r)
	location := w.HeaderMap["Location"]
	if len(location) == 0 {
		t.Fatal("No 'Location' header empty or missed in the redirect response\n")
	}
	got := location[0]
	want := "http://google.com"

	if got != want {
		t.Fatalf("expected to have %q location after redirect, but got %q\n", want, got)
	}

}
