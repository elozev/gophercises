package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	urlshort "urlshort/src"

	bolt "go.etcd.io/bbolt"
)

var (
	readFile  = flag.String("file", "", "Specify a yaml or json file to read from")
	dbEnabled = flag.Bool("db", false, "Use BoltDB")
)

func initDb() *bolt.DB {
	db, err := bolt.Open("redirects.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to DB!")

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Redirects"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Redirects"))
		err = bucket.Put([]byte("/dog"), []byte("https://http.dog"))
		return err
	})

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Redirects"))
		err = bucket.Put([]byte("/urlshort"), []byte("https://github.com/gophercises/urlshort"))
		return err
	})

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Redirects"))
		err = bucket.Put([]byte("/urlshort-final"), []byte("https://github.com/gophercises/urlshort/tree/solution"))
		return err
	})

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Redirects"))
		err = bucket.Put([]byte("/cat"), []byte("https://http.cat"))
		return err
	})

	return db
}

func main() {
	var db *bolt.DB

	mux := defaultMux()
	flag.Parse()

	if *dbEnabled {
		db = initDb()
		defer db.Close()
	}

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

	var handler http.Handler

	if *dbEnabled {
		handler, err = urlshort.DbHandler(db, mux)
		if err != nil {
			panic(err)
		}
	} else {
		handler, err = urlshort.Handler([]byte(data), dataType, mux)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
