package url

import "net/http"

type RedirectHandler struct{}

func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://google.com", http.StatusFound)
}
