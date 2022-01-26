package urlshort

import (
	"net/http"
)

type RedirectHandler struct {
	redirects map[string]string
	fallback  http.Handler
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

func MapHandler(redirects map[string]string, fallback http.Handler) RedirectHandler {
	return RedirectHandler{redirects, fallback}
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

func parseYAML(yaml []byte) ([]redirect, error) {
	res := make([]redirect, 0)
	return res, nil
}
