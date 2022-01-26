package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
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

func MapHandler(redirects map[string]string, fallback http.Handler) http.Handler {
	return RedirectHandler{redirects, fallback}
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.Handler, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return RedirectHandler{}, nil
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
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

func parseYAML(yamlBinary []byte) ([]redirect, error) {
	res := make([]redirect, 0)
	err := yaml.Unmarshal(yamlBinary, &res)
	return res, err
}
