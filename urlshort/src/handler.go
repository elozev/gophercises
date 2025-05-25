package urlshort

import (
	"fmt"
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"

	"gopkg.in/yaml.v3"
)

const (
	DATA_TYPE_JSON = "json"
	DATA_TYPE_YAML = "yaml"
	DATA_TYPE_DB   = "db"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.RequestURI

		if redirect, ok := pathsToUrls[url]; ok {
			http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

type Redirect struct {
	From string `yaml:"path" json:"path"`
	To   string `yaml:"url" json:"url"`
}

func parseYaml(yml []byte) ([]Redirect, error) {
	var redirects []Redirect
	err := yaml.Unmarshal(yml, &redirects)
	if err != nil {
		return nil, err
	}
	return redirects, nil
}

func parseJson(json []byte) ([]Redirect, error) {
	var redirects []Redirect
	err := yaml.Unmarshal(json, &redirects)
	if err != nil {
		return nil, err
	}

	return redirects, nil
}

func mapBuilder(redirects []Redirect) map[string]string {
	res := make(map[string]string, len(redirects))

	for _, r := range redirects {
		res[r.From] = r.To
	}

	return res
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
	redirects, err := parseYaml(yml)

	if err != nil {
		return nil, err
	}

	pathsToUrls := mapBuilder(redirects)
	return MapHandler(pathsToUrls, fallback), nil
}

func JsonHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	redirects, err := parseJson(json)

	if err != nil {
		return nil, err
	}

	pathsToUrls := mapBuilder(redirects)
	return MapHandler(pathsToUrls, fallback), nil
}

func Handler(data []byte, dataType string, fallback http.Handler) (http.HandlerFunc, error) {
	if dataType == DATA_TYPE_JSON {
		return JsonHandler(data, fallback)
	} else if dataType == DATA_TYPE_YAML {
		return YAMLHandler(data, fallback)
	} else {
		return nil, fmt.Errorf("data type \"%s\" not supported", dataType)
	}
}

func DbHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.RequestURI

		var value []byte

		err := db.View(func(tx *bolt.Tx) error {
			log.Println("viewing")
			b := tx.Bucket([]byte("Redirects"))
			if b == nil {
				return fmt.Errorf("Redirects bucket not found")
			}

			value = b.Get([]byte(url))
			if value == nil {
				return fmt.Errorf("key not found")
			}
			log.Printf("value for %s: %s", url, string(value))

			return nil
		})

		if err != nil {
			log.Printf("Error: %v", err)
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, string(value), http.StatusPermanentRedirect)
	}, nil
}
