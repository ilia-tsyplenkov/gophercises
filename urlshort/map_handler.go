package short

import (
	"net/http"
)

type RedirectHandler struct {
	Redirects map[string]string
}

func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, value := range h.Redirects {
		if key == r.URL.Path {
			http.Redirect(w, r, value, http.StatusFound)
			return

		}
	}
}
