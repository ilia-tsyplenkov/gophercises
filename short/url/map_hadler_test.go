package url_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilia-tsyplenkov/gophercises/short/url"
)

func TestRedirectHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	// r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler := url.RedirectHandler{}
	handler.ServeHTTP(w, r)
	want := http.StatusFound
	got := w.Code
	if got != want {
		t.Fatalf("expected to have %d code but got %d\n", want, got)
	}
}
