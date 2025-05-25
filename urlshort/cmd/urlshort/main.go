package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	urlshort "urlshort/src"
)

var (
	yamlFile = flag.String("yaml-file", "", "Specify a yaml file to read from")
)

func main() {
	mux := defaultMux()
	flag.Parse()

	var yaml []byte
	var err error

	if *yamlFile != "" {
		log.Println("Yaml file specified, reading!")
		yaml, err = urlshort.ParseFile(*yamlFile)

		if err != nil {
			log.Fatalln(err)
		}

	} else {
	// Build the YAMLHandler using the mapHandler as the
	// fallback
		yaml = []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mux)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}