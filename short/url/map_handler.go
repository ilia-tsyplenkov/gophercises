package url

import "net/http"

type RedirectHandler struct{}

func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}
