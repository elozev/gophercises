package main

import (
	cyoa "cyoa/src"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	filename = flag.String("filename", "gopher.json", "Filename to be used as a base for the story program!")
	port     = flag.Int("port", 8080, "Default port to run the program on")
)

func main() {
	flag.Parse()
	fmt.Println("Web program!")

	var sh cyoa.StoryHandler

	// Open the gopher.json file
	file, err := os.Open(*filename)
	cyoa.Check(err)
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	cyoa.Check(err)

	var storiesHolder cyoa.Story

	err = json.Unmarshal(fileContents, &storiesHolder)
	cyoa.Check(err)

	sh.Stories = storiesHolder

	mux := http.NewServeMux()

	mux.HandleFunc("/", sh.IndexHandler)
	mux.HandleFunc("/story/{chapter}", sh.StoryHandler)

	log.Printf("Listening and serving on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
