package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	decodedYaml := []map[string]string{}
	err := yaml.Unmarshal(yml, &decodedYaml)

	pathsToUrls := createMap(decodedYaml)

	return MapHandler(pathsToUrls, fallback), err
}

// CreateMap uses the provided map slice, creates a new map and for each map
// present in the provided slice , it creates a entry in the newly created map by
// using path as key and url as key value.
func createMap(decodedYaml []map[string]string) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, yamlMap := range decodedYaml {
		key := yamlMap["path"]
		pathsToUrls[key] = yamlMap["url"]
	}
	return pathsToUrls
}
