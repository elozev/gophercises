package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	urlshort "urlshort/src"
)

var (
	readFile = flag.String("file", "", "Specify a yaml or json file to read from")
)

func main() {
	mux := defaultMux()
	flag.Parse()

	var data []byte
	var err error
	var dataType = urlshort.DATA_TYPE_YAML

	if *readFile != "" {
		log.Println("File specified, reading!")
		data, err = urlshort.ParseFile(*readFile)

		if err != nil {
			log.Fatalln(err)
		}

		fileExtension := filepath.Ext(*readFile) 

		if fileExtension == "yaml" || fileExtension == "yml" {
			dataType = urlshort.DATA_TYPE_YAML
		} else if fileExtension == "json" {
			dataType = urlshort.DATA_TYPE_JSON
		}
		
	} else {
	// Build the YAMLHandler using the mapHandler as the
	// fallback
		data = []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	}

	yamlHandler, err := urlshort.Handler([]byte(data),dataType, mux)
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