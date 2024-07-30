package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// return func with hanlder response writer and request
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusOK)
		}
		fallback.ServeHTTP(w, r)
	}
	// look if the path is present in our map
	// if yes redirect it
	// else use fallback

}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	// 2 parse the yaml
	var pathUrls []pathUrl

	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}

	// convert it to a map
	pathsToUrls := make(map[string]string)

	for _, value := range pathUrls {
		pathsToUrls[value.PATH] = value.URL
	}
	// return the map handler

	return MapHandler(pathsToUrls, fallback), nil
}

// 1 initialize the struct

type pathUrl struct {
	PATH string `json:"path,omitempty"`
	URL  string `json:"url,omitempty"`
}
