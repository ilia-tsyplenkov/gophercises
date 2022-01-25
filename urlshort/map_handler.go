package short

import (
	"net/http"
)

type RedirectHandler struct {
	redirects map[string]string
	fallback  http.Handler
}

func NewRedirectHandler(redirects map[string]string, fallback http.Handler) RedirectHandler {
	return RedirectHandler{redirects, fallback}
}
func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, value := range h.redirects {
		if key == r.URL.Path {
			http.Redirect(w, r, value, http.StatusFound)
			return

		}
	}
	h.fallback.ServeHTTP(w, r)
}
